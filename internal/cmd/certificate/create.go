package certificate

import (
	"context"
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
			Use:   "create [FLAGS]",
			Short: "Create or upload a Certificate",
			Args:  cobra.ExactArgs(0),
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
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, strings []string) (*hcloud.Response, any, error) {
		certType, err := cmd.Flags().GetString("type")
		if err != nil {
			return nil, nil, err
		}
		switch hcloud.CertificateType(certType) {
		case hcloud.CertificateTypeManaged:
			response, err := createManaged(ctx, client, waiter, cmd)
			return response, nil, err
		default: // Uploaded
			response, err := createUploaded(ctx, client, cmd)
			return response, nil, err
		}
	},
	PrintResource: func(_ context.Context, _ hcapi2.Client, _ *cobra.Command, _ any) {
		// no-op
	},
}

func createUploaded(ctx context.Context, client hcapi2.Client, cmd *cobra.Command) (*hcloud.Response, error) {
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
	cert, response, err := client.Certificate().Create(ctx, createOpts)
	if err != nil {
		return nil, err
	}
	cmd.Printf("Certificate %d created\n", cert.ID)
	return response, nil
}

func createManaged(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command) (*hcloud.Response, error) {
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
	res, response, err := client.Certificate().CreateCertificate(ctx, createOpts)
	if err != nil {
		return nil, err
	}
	if err := waiter.ActionProgress(ctx, res.Action); err != nil {
		return nil, err
	}
	cmd.Printf("Certificate %d created\n", res.Certificate.ID)
	return response, nil
}
