package certificate

import (
	"context"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

var listCmd = base.ListCmd{
	ResourceNamePlural: "certificates",
	DefaultColumns:     []string{"id", "name", "type", "domain_names", "not_valid_after"},

	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, listOpts hcloud.ListOpts, sorts []string) ([]interface{}, error) {
		opts := hcloud.CertificateListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		certificates, err := client.Certificate().AllWithOpts(ctx, opts)

		var resources []interface{}
		for _, n := range certificates {
			resources = append(resources, n)
		}
		return resources, err
	},

	OutputTable: func(_ hcapi2.Client) *output.Table {
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
			AddFieldFn("issuance_status", func(obj interface{}) string {
				cert := obj.(*hcloud.Certificate)
				if cert.Type != hcloud.CertificateTypeManaged {
					return "n/a"
				}
				return string(cert.Status.Issuance)
			}).
			AddFieldFn("renewal_status", func(obj interface{}) string {
				cert := obj.(*hcloud.Certificate)
				if cert.Type != hcloud.CertificateTypeManaged ||
					cert.Status.Renewal == hcloud.CertificateStatusTypeUnavailable {
					return "n/a"
				}
				return string(cert.Status.Renewal)
			}).
			AddFieldFn("domain_names", func(obj interface{}) string {
				cert := obj.(*hcloud.Certificate)
				return strings.Join(cert.DomainNames, ", ")
			}).
			AddFieldFn("created", output.FieldFn(func(obj interface{}) string {
				cert := obj.(*hcloud.Certificate)
				return util.Datetime(cert.Created)
			}))
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var certSchemas []schema.Certificate
		for _, resource := range resources {
			cert := resource.(*hcloud.Certificate)
			certSchema := schema.Certificate{
				ID:             cert.ID,
				Certificate:    cert.Certificate,
				Created:        cert.Created,
				DomainNames:    cert.DomainNames,
				Fingerprint:    cert.Fingerprint,
				Labels:         cert.Labels,
				Name:           cert.Name,
				Type:           string(cert.Type),
				NotValidAfter:  cert.NotValidAfter,
				NotValidBefore: cert.NotValidBefore,
			}
			if len(cert.UsedBy) > 0 {
				certSchema.UsedBy = make([]schema.CertificateUsedByRef, len(cert.UsedBy))
				for i, ub := range cert.UsedBy {
					certSchema.UsedBy[i] = schema.CertificateUsedByRef{
						ID:   ub.ID,
						Type: string(ub.Type),
					}
				}
			}
			certSchemas = append(certSchemas, certSchema)
		}

		return certSchemas
	},
}
