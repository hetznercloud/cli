package certificate

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var CreateCmd = base.CreateCmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create [options] --name <name> (--type managed --domain <domain> | --type uploaded --cert-file <file> --key-file <file>)",
			Short: "Create or upload a Certificate",
		}

		cmd.Flags().String("name", "", "Certificate name (required)")
		cmd.MarkFlagRequired("name")

		cmd.Flags().StringP("type", "t", string(hcloud.CertificateTypeUploaded),
			fmt.Sprintf("Type of certificate to create. Valid choices: %v, %v",
				hcloud.CertificateTypeUploaded, hcloud.CertificateTypeManaged))
		cmd.RegisterFlagCompletionFunc(
			"type",
			cmpl.SuggestCandidates(string(hcloud.CertificateTypeUploaded), string(hcloud.CertificateTypeManaged)),
		)

		cmd.Flags().String("cert-file", "", "File containing the PEM encoded certificate (required if type is uploaded)")
		cmd.Flags().String("key-file", "",
			"File containing the PEM encoded private key for the certificate (required if type is uploaded)")
		cmd.Flags().StringSlice("domain", nil, "One or more domains the certificate is valid for.")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, strings []string) (any, any, error) {
		certType, err := cmd.Flags().GetString("type")
		if err != nil {
			return nil, nil, err
		}
		var cert *hcloud.Certificate
		switch hcloud.CertificateType(certType) {
		case hcloud.CertificateTypeManaged:
			cert, err = createManaged(s, cmd)
		default: // Uploaded
			cert, err = createUploaded(s, cmd)
		}
		if err != nil {
			return nil, nil, err
		}
		return cert, util.Wrap("certificate", hcloud.SchemaFromCertificate(cert)), nil
	},
}

func createUploaded(s state.State, cmd *cobra.Command) (*hcloud.Certificate, error) {
	var (
		name              string
		certFile, keyFile string
		certPEM, keyPEM   []byte
		cert              *hcloud.Certificate

		err error
	)

	if err = util.ValidateRequiredFlags(cmd.Flags(), "cert-file", "key-file"); err != nil {
		return nil, err
	}
	if name, err = cmd.Flags().GetString("name"); err != nil {
		return nil, err
	}
	if certFile, err = cmd.Flags().GetString("cert-file"); err != nil {
		return nil, err
	}
	if keyFile, err = cmd.Flags().GetString("key-file"); err != nil {
		return nil, err
	}

	if certPEM, err = os.ReadFile(certFile); err != nil {
		return nil, err
	}
	if keyPEM, err = os.ReadFile(keyFile); err != nil {
		return nil, err
	}

	createOpts := hcloud.CertificateCreateOpts{
		Name:        name,
		Type:        hcloud.CertificateTypeUploaded,
		Certificate: string(certPEM),
		PrivateKey:  string(keyPEM),
	}
	cert, _, err = s.Client().Certificate().Create(s, createOpts)
	if err != nil {
		return nil, err
	}
	cmd.Printf("Certificate %d created\n", cert.ID)
	return cert, nil
}

func createManaged(s state.State, cmd *cobra.Command) (*hcloud.Certificate, error) {
	var (
		name    string
		domains []string
		res     hcloud.CertificateCreateResult
		err     error
	)

	if name, err = cmd.Flags().GetString("name"); err != nil {
		return nil, nil
	}
	if err = util.ValidateRequiredFlags(cmd.Flags(), "domain"); err != nil {
		return nil, err
	}
	if domains, err = cmd.Flags().GetStringSlice("domain"); err != nil {
		return nil, nil
	}

	createOpts := hcloud.CertificateCreateOpts{
		Name:        name,
		Type:        hcloud.CertificateTypeManaged,
		DomainNames: domains,
	}
	res, _, err = s.Client().Certificate().CreateCertificate(s, createOpts)
	if err != nil {
		return nil, err
	}
	if err := s.WaitForActions(cmd, s, res.Action); err != nil {
		return nil, err
	}
	defer cmd.Printf("Certificate %d created\n", res.Certificate.ID)
	cert, _, err := s.Client().Certificate().GetByID(s, res.Certificate.ID)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
