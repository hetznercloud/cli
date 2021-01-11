package cmds

import (
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var certificateTableOutput *tableOutput

func init() {
	certificateTableOutput = describeCertificatesTableOutput()
}

func newCertificatesListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Certificates",
		Long: util.ListLongDescription(
			"Displays a list of certificates",
			certificateTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runCertificatesList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(serverListTableOutput.Columns()), outputOptionJSON())
	return cmd
}

func runCertificatesList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")

	opts := hcloud.CertificateListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}

	certs, err := cli.Client().Certificate.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var certSchemas []schema.Certificate

		for _, cert := range certs {
			certSchema := schema.Certificate{
				ID:             cert.ID,
				Certificate:    cert.Certificate,
				Created:        cert.Created,
				DomainNames:    cert.DomainNames,
				Fingerprint:    cert.Fingerprint,
				Labels:         cert.Labels,
				Name:           cert.Name,
				NotValidAfter:  cert.NotValidAfter,
				NotValidBefore: cert.NotValidBefore,
			}
			certSchemas = append(certSchemas, certSchema)
		}

		return util.DescribeJSON(certSchemas)
	}

	cols := []string{"id", "name", "domain_names", "not_valid_after"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}
	tw := describeCertificatesTableOutput()
	if err := tw.ValidateColumns(cols); err != nil {
		return nil
	}
	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, cert := range certs {
		tw.Write(cols, cert)
	}
	return tw.Flush()
}

func describeCertificatesTableOutput() *tableOutput {
	return newTableOutput().
		AddAllowedFields(hcloud.Certificate{}).
		RemoveAllowedField("certificate", "chain").
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.LabelsToString(cert.Labels)
		})).
		AddFieldOutputFn("not_valid_before", func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.Datetime(cert.NotValidBefore)
		}).
		AddFieldOutputFn("not_valid_after", func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.Datetime(cert.NotValidAfter)
		}).
		AddFieldOutputFn("domain_names", func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return strings.Join(cert.DomainNames, ", ")
		}).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.Datetime(cert.Created)
		}))
}
