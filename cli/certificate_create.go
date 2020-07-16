package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newCertificateCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS]",
		Short:                 "Create or upload a Certificate",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runCertificateCreate),
	}

	cmd.Flags().String("name", "", "Certificate name (required)")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("cert-file", "", "File containing the PEM encoded certificate (required)")
	cmd.MarkFlagRequired("cert-file")

	cmd.Flags().String("key-file", "", "File containing the PEM encoded private key for the certificate (required)")
	cmd.MarkFlagRequired("key-file")

	return cmd
}

func runCertificateCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	var (
		name string

		certFile, keyFile string
		certPEM, keyPEM   []byte
		cert              *hcloud.Certificate

		err error
	)
	if name, err = cmd.Flags().GetString("name"); err != nil {
		return err
	}
	if certFile, err = cmd.Flags().GetString("cert-file"); err != nil {
		return err
	}
	if keyFile, err = cmd.Flags().GetString("key-file"); err != nil {
		return err
	}

	if certPEM, err = ioutil.ReadFile(certFile); err != nil {
		return err
	}
	if keyPEM, err = ioutil.ReadFile(keyFile); err != nil {
		return err
	}

	createOpts := hcloud.CertificateCreateOpts{
		Certificate: string(certPEM),
		Name:        name,
		PrivateKey:  string(keyPEM),
	}
	if cert, _, err = cli.Client().Certificate.Create(cli.Context, createOpts); err != nil {
		return err
	}
	fmt.Printf("Certificate %d created\n", cert.ID)
	return nil
}
