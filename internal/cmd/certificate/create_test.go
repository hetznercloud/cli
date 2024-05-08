package certificate_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/certificate"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/managed_create_response.json
var managedCreateResponseJson string

//go:embed testdata/uploaded_create_response.json
var uploadedCreateResponseJson string

func TestCreateManaged(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		CreateCertificate(gomock.Any(), hcloud.CertificateCreateOpts{
			Name:        "test",
			Type:        hcloud.CertificateTypeManaged,
			DomainNames: []string{"example.com"},
		}).
		Return(hcloud.CertificateCreateResult{
			Certificate: &hcloud.Certificate{
				ID:          123,
				Name:        "test",
				Type:        hcloud.CertificateTypeManaged,
				DomainNames: []string{"example.com"},
			},
			Action: &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321})
	fx.Client.CertificateClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&hcloud.Certificate{
			ID:          123,
			Name:        "test",
			Type:        hcloud.CertificateTypeManaged,
			DomainNames: []string{"example.com"},
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"--name", "test", "--type", "managed", "--domain", "example.com"})

	expOut := "Certificate 123 created\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateManagedJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		CreateCertificate(gomock.Any(), hcloud.CertificateCreateOpts{
			Name:        "test",
			Type:        hcloud.CertificateTypeManaged,
			DomainNames: []string{"example.com"},
		}).
		Return(hcloud.CertificateCreateResult{
			Certificate: &hcloud.Certificate{
				ID:             123,
				Name:           "test",
				Type:           hcloud.CertificateTypeManaged,
				Created:        time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
				NotValidBefore: time.Time{},
				NotValidAfter:  time.Time{},
				DomainNames:    []string{"example.com"},
				Labels:         map[string]string{"key": "value"},
				UsedBy: []hcloud.CertificateUsedByRef{{
					ID:   123,
					Type: hcloud.CertificateUsedByRefTypeLoadBalancer,
				}},
				Status: &hcloud.CertificateStatus{
					Issuance: hcloud.CertificateStatusTypePending,
					Renewal:  hcloud.CertificateStatusTypeUnavailable,
				},
			},
			Action: &hcloud.Action{ID: 321},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321})
	fx.Client.CertificateClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&hcloud.Certificate{
			ID:             123,
			Name:           "test",
			Type:           hcloud.CertificateTypeManaged,
			Created:        time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
			NotValidBefore: time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
			NotValidAfter:  time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
			DomainNames:    []string{"example.com"},
			Labels:         map[string]string{"key": "value"},
			UsedBy: []hcloud.CertificateUsedByRef{{
				ID:   123,
				Type: hcloud.CertificateUsedByRefTypeLoadBalancer,
			}},
			Status: &hcloud.CertificateStatus{
				Issuance: hcloud.CertificateStatusTypeCompleted,
				Renewal:  hcloud.CertificateStatusTypeUnavailable,
			},
			Fingerprint: "fingerprint placeholder",
			Certificate: "certificate data placeholder",
		}, nil, nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "test", "--type", "managed", "--domain", "example.com"})

	expOut := "Certificate 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)

	assert.JSONEq(t, managedCreateResponseJson, jsonOut)
}

func TestCreateUploaded(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		Create(gomock.Any(), hcloud.CertificateCreateOpts{
			Name:        "test",
			Type:        hcloud.CertificateTypeUploaded,
			Certificate: "certificate file content",
			PrivateKey:  "key file content",
		}).
		Return(&hcloud.Certificate{
			ID:   123,
			Name: "test",
			Type: hcloud.CertificateTypeUploaded,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"--name", "test", "--key-file", "testdata/key.pem", "--cert-file", "testdata/cert.pem"})

	expOut := "Certificate 123 created\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestCreateUploadedJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := certificate.CreateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.CertificateClient.EXPECT().
		Create(gomock.Any(), hcloud.CertificateCreateOpts{
			Name:        "test",
			Type:        hcloud.CertificateTypeUploaded,
			Certificate: "certificate file content",
			PrivateKey:  "key file content",
		}).
		Return(&hcloud.Certificate{
			ID:             123,
			Name:           "test",
			Type:           hcloud.CertificateTypeUploaded,
			Created:        time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
			NotValidBefore: time.Date(2020, 8, 24, 12, 0, 0, 0, time.UTC),
			NotValidAfter:  time.Date(2036, 8, 12, 12, 0, 0, 0, time.UTC),
			Labels:         map[string]string{"key": "value"},
			Fingerprint:    "00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00",
			UsedBy: []hcloud.CertificateUsedByRef{{
				ID:   123,
				Type: hcloud.CertificateUsedByRefTypeLoadBalancer,
			}},
		}, nil, nil)

	jsonOut, out, err := fx.Run(cmd, []string{"-o=json", "--name", "test", "--key-file", "testdata/key.pem", "--cert-file", "testdata/cert.pem"})

	expOut := "Certificate 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.JSONEq(t, uploadedCreateResponseJson, jsonOut)
}
