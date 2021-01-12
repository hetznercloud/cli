package certificate

import (
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var listTableOutput *output.Table

func init() {
	listTableOutput = describeTableOutput()
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Certificates",
		Long: util.ListLongDescription(
			"Displays a list of certificates",
			listTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(listTableOutput.Columns()), output.OptionJSON())
	return cmd
}

func runList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

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
	tw := describeTableOutput()
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

func describeTableOutput() *output.Table {
	return output.NewTable().
		AddAllowedFields(hcloud.Certificate{}).
		RemoveAllowedField("certificate", "chain").
		AddFieldFn("labels", output.FieldFn(func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.LabelsToString(cert.Labels)
		})).
		AddFieldFn("not_valid_before", func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.Datetime(cert.NotValidBefore)
		}).
		AddFieldFn("not_valid_after", func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.Datetime(cert.NotValidAfter)
		}).
		AddFieldFn("domain_names", func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return strings.Join(cert.DomainNames, ", ")
		}).
		AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
			cert := obj.(*hcloud.Certificate)
			return util.Datetime(cert.Created)
		}))
}
