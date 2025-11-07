//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/kit/randutil"
)

const fingerprintRegex = `[0-9A-F]{2}(:[0-9A-F]{2}){31}`

func TestCertificate(t *testing.T) {
	t.Parallel()

	t.Run("uploaded", func(t *testing.T) {
		tmpDir := t.TempDir()
		notBefore := time.Now()
		notAfter := notBefore.Add(365 * 24 * time.Hour)
		certPath, keyPath := path.Join(tmpDir, "cert.pem"), path.Join(tmpDir, "key.pem")
		err := generateCertificate(certPath, keyPath, notBefore, notAfter)
		require.NoError(t, err)

		certName := withSuffix("test-certificate-uploaded")
		certID, err := createCertificate(t, certName, hcloud.CertificateTypeUploaded, "--cert-file", certPath, "--key-file", keyPath)
		require.NoError(t, err)

		runCertificateTestSuite(t, certName, certID, hcloud.CertificateTypeUploaded, "example.com")
	})

	t.Run("managed", func(t *testing.T) {
		certDomain := os.Getenv("CERT_DOMAIN")
		if certDomain == "" {
			t.Skip("Skipping because CERT_DOMAIN is not set")
		}

		// random subdomain
		certDomain = fmt.Sprintf("%s.%s", randutil.GenerateID(), certDomain)

		certName := withSuffix("test-certificate-managed")
		certID, err := createCertificate(t, certName, hcloud.CertificateTypeManaged, "--type", "managed", "--domain", certDomain)
		require.NoError(t, err)

		runCertificateTestSuite(t, certName, certID, hcloud.CertificateTypeManaged, certDomain)
	})
}

func runCertificateTestSuite(t *testing.T, certName string, certID int64, certType hcloud.CertificateType, domainName string) {
	t.Helper()

	t.Run("add-label", func(t *testing.T) {
		t.Run("non-existing", func(t *testing.T) {
			out, err := runCommand(t, "certificate", "add-label", "non-existing-certificate", "foo=bar")
			require.EqualError(t, err, "Certificate not found: non-existing-certificate")
			assert.Empty(t, out)
		})

		t.Run("1", func(t *testing.T) {
			out, err := runCommand(t, "certificate", "add-label", strconv.FormatInt(certID, 10), "foo=bar")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) foo added to Certificate %d\n", certID), out)
		})

		t.Run("2", func(t *testing.T) {
			out, err := runCommand(t, "certificate", "add-label", strconv.FormatInt(certID, 10), "baz=qux")
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("Label(s) baz added to Certificate %d\n", certID), out)
		})
	})

	t.Run("update-name", func(t *testing.T) {
		certName = withSuffix("new-test-certificate-" + string(certType))

		out, err := runCommand(t, "certificate", "update", strconv.FormatInt(certID, 10), "--name", certName)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Certificate %d updated\n", certID), out)
	})

	t.Run("remove-label", func(t *testing.T) {
		out, err := runCommand(t, "certificate", "remove-label", certName, "baz")
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Label(s) baz removed from Certificate %d\n", certID), out)
	})

	t.Run("list", func(t *testing.T) {
		out, err := runCommand(t, "certificate", "list", "-o=columns=id,name,labels,type,created,"+
			"not_valid_before,not_valid_after,domain_names,fingerprint,issuance_status,renewal_status,age")
		require.NoError(t, err)

		labels := []string{"foo=bar"}
		if certType == hcloud.CertificateTypeManaged {
			labels = append([]string{"HC-Use-Staging-CA=true"}, labels...)
		}

		assert.Regexp(t,
			NewRegex().Start().
				SeparatedByWhitespace(
					"ID", "NAME", "LABELS", "TYPE", "CREATED", "NOT VALID BEFORE", "NOT VALID AFTER",
					"DOMAIN NAMES", "FINGERPRINT", "ISSUANCE STATUS", "RENEWAL STATUS", "AGE",
				).OptionalWhitespace().Newline().
				Lit(strconv.FormatInt(certID, 10)).Whitespace().
				Lit(certName).Whitespace().
				Lit(strings.Join(labels, ", ")).Whitespace().
				Lit(string(certType)).Whitespace().
				Datetime().Whitespace().
				Datetime().Whitespace().
				Datetime().Whitespace().
				Lit(domainName).Whitespace().
				Raw(fingerprintRegex).Whitespace().
				OneOf("completed", "n/a").Whitespace().
				Lit("n/a").Whitespace().
				Age().OptionalWhitespace().Newline().
				End(),
			out,
		)
	})

	t.Run("describe", func(t *testing.T) {
		out, err := runCommand(t, "certificate", "describe", strconv.FormatInt(certID, 10))
		require.NoError(t, err)

		regex := NewRegex().Start().
			Lit("ID:").Whitespace().Int().Newline().
			Lit("Name:").Whitespace().Lit(certName).Newline().
			Lit("Type:").Whitespace().Lit(string(certType)).Newline().
			Lit("Fingerprint:").Whitespace().Raw(fingerprintRegex).Newline().
			Lit("Created:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
			Lit("Not valid before:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
			Lit("Not valid after:").Whitespace().Datetime().Lit(" (").HumanizeTime().Lit(")").Newline().
			Newline()

		if certType == hcloud.CertificateTypeManaged {
			regex = regex.
				Lit("Status:").Newline().
				Lit("  Issuance:").Whitespace().Lit("completed").Newline().
				Lit("  Renewal:").Whitespace().Lit("unavailable").Newline().
				Newline()
		}

		regex = regex.
			Lit("Domain names:").Newline().
			Lit("  - ").Lit(domainName).Newline().
			Newline().
			Lit("Labels:").Newline()

		if certType == hcloud.CertificateTypeManaged {
			regex = regex.Lit("  HC-Use-Staging-CA:").Whitespace().Lit("true").Newline()
		}

		regex = regex.
			Lit("  foo:").Whitespace().Lit("bar").Newline().
			Newline().
			Lit("Used By:").Newline().
			Lit("  Certificate unused").Newline().
			End()

		assert.Regexp(t, regex, out)
	})

	t.Run("retry", func(t *testing.T) {
		out, err := runCommand(t, "certificate", "retry", strconv.FormatInt(certID, 10))
		assert.Empty(t, out)
		require.Error(t, err)
		assert.Regexp(t, `certificate not retryable \(unsupported_error, [0-9a-f]+\)`, err.Error())
	})

	t.Run("delete", func(t *testing.T) {
		out, err := runCommand(t, "certificate", "delete", strconv.FormatInt(certID, 10))
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("Certificate %d deleted\n", certID), out)
	})
}

func createCertificate(t *testing.T, name string, certificateType hcloud.CertificateType, args ...string) (int64, error) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = client.Certificate.Delete(context.Background(), &hcloud.Certificate{Name: name})
	})

	if certificateType == hcloud.CertificateTypeManaged {
		args = append([]string{"--label", "HC-Use-Staging-CA=true"}, args...)
	}

	out, err := runCommand(t, append([]string{"certificate", "create", "--name", name, "--type", string(certificateType)}, args...)...)
	if err != nil {
		return 0, err
	}

	if !assert.Regexp(t, `^Certificate [0-9]+ created\n$`, out) {
		return 0, fmt.Errorf("invalid response: %s", out)
	}

	id, err := strconv.ParseInt(out[12:len(out)-9], 10, 64)
	if err != nil {
		return 0, err
	}

	t.Cleanup(func() {
		_, _ = client.Certificate.Delete(context.Background(), &hcloud.Certificate{ID: id})
	})
	return id, nil
}
