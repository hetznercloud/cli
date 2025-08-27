package certificate

import (
	"strings"
	"time"

	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var ListCmd = &base.ListCmd[*hcloud.Certificate, schema.Certificate]{
	ResourceNamePlural: "Certificates",
	JSONKeyGetByName:   "certificates",
	DefaultColumns:     []string{"id", "name", "type", "domain_names", "not_valid_after", "age"},
	SortOption:         config.OptionSortCertificate,

	Fetch: func(s state.State, _ *pflag.FlagSet, listOpts hcloud.ListOpts, sorts []string) ([]*hcloud.Certificate, error) {
		opts := hcloud.CertificateListOpts{ListOpts: listOpts}
		if len(sorts) > 0 {
			opts.Sort = sorts
		}
		return s.Client().Certificate().AllWithOpts(s, opts)
	},

	OutputTable: func(t *output.Table, _ hcapi2.Client) {
		t.
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
			})).
			AddFieldFn("age", output.FieldFn(func(obj interface{}) string {
				cert := obj.(*hcloud.Certificate)
				return util.Age(cert.Created, time.Now())
			}))
	},

	Schema: hcloud.SchemaFromCertificate,
}
