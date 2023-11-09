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

var CreateCmd = base.Cmd{
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
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, strings []string) error {
		certType, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}
		switch hcloud.CertificateType(certType) {
		case hcloud.CertificateTypeUploaded:
			return createUploaded(ctx, client, cmd)
		case hcloud.CertificateTypeManaged:
			return createManaged(ctx, client, waiter, cmd)
		default:
			return createUploaded(ctx, client, cmd)
		}
	},
}

func createUploaded(ctx context.Context, client hcapi2.Client, cmd *cobra.Command) error {
	var (
		name              string
		certFile, keyFile string
		certPEM, keyPEM   []byte
		cert              *hcloud.Certificate

		err error
	)

	if err = util.ValidateRequiredFlags(cmd.Flags(), "cert-file", "key-file"); err != nil {
		return err
	}
	if name, err = cmd.Flags().GetString("name"); err != nil {
		return err
	}
	if certFile, err = cmd.Flags().GetString("cert-file"); err != nil {
		return err
	}
	if keyFile, err = cmd.Flags().GetString("key-file"); err != nil {
		return err
	}

	if certPEM, err = os.ReadFile(certFile); err != nil {
		return err
	}
	if keyPEM, err = os.ReadFile(keyFile); err != nil {
		return err
	}

	createOpts := hcloud.CertificateCreateOpts{
		Name:        name,
		Type:        hcloud.CertificateTypeUploaded,
		Certificate: string(certPEM),
		PrivateKey:  string(keyPEM),
	}
	if cert, _, err = client.Certificate().Create(ctx, createOpts); err != nil {
		return err
	}
	cmd.Printf("Certificate %d created\n", cert.ID)
	return nil
}

func createManaged(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command) error {
	var (
		name    string
		domains []string
		res     hcloud.CertificateCreateResult
		err     error
	)

	if name, err = cmd.Flags().GetString("name"); err != nil {
		return nil
	}
	if err = util.ValidateRequiredFlags(cmd.Flags(), "domain"); err != nil {
		return err
	}
	if domains, err = cmd.Flags().GetStringSlice("domain"); err != nil {
		return nil
	}

	createOpts := hcloud.CertificateCreateOpts{
		Name:        name,
		Type:        hcloud.CertificateTypeManaged,
		DomainNames: domains,
	}
	if res, _, err = client.Certificate().CreateCertificate(ctx, createOpts); err != nil {
		return err
	}
	if err := waiter.ActionProgress(ctx, res.Action); err != nil {
		return err
	}
	cmd.Printf("Certificate %d created\n", res.Certificate.ID)
	return nil
}
