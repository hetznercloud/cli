package util

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func YesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func NA(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func Datetime(t time.Time) string {
	return t.Local().Format(time.UnixDate)
}

func Age(t time.Time) string {
	currentTime := time.Now()
	diff := currentTime.Sub(t)

	if diff.Hours() >= 24 {
		days := int(diff.Hours()) / 24
		return fmt.Sprintf("%dd", days)
	}

	if diff.Hours() > 0 {
		return fmt.Sprintf("%dh", int(diff.Hours()))
	}

	if diff.Minutes() > 0 {
		return fmt.Sprintf("%dm", int(diff.Minutes()))
	}

	if diff.Seconds() > 0 {
		return fmt.Sprintf("%ds", int(diff.Seconds()))
	}

	return "just now"
}

func ChainRunE(fns ...func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, fn := range fns {
			if fn == nil {
				continue
			}
			if err := fn(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

func ExactlyOneSet(s string, ss ...string) bool {
	set := s != ""
	for _, s := range ss {
		if set && s != "" {
			return false
		}
		set = set || s != ""
	}
	return set
}

var outputDescription = `Output can be controlled with the -o flag. Use -o noheader to suppress the
table header. Displayed columns and their order can be set with
-o columns=%s (see available columns below).`

func ListLongDescription(intro string, columns []string) string {
	var colExample []string
	if len(columns) > 2 {
		colExample = columns[0:2]
	} else {
		colExample = columns
	}
	return fmt.Sprintf(
		"%s\n\n%s\n\nColumns:\n - %s",
		intro,
		fmt.Sprintf(outputDescription, strings.Join(colExample, ",")),
		strings.Join(columns, "\n - "),
	)
}

func SplitLabel(label string) []string {
	return strings.SplitN(label, "=", 2)
}

// SplitLabelVars splits up label into key and value and returns them as separate return values.
// If label doesn't contain the `=` separator, SplitLabelVars returns the original string as key,
// with an empty value.
func SplitLabelVars(label string) (string, string) {
	parts := strings.SplitN(label, "=", 2)
	if len(parts) != 2 {
		return label, ""
	}
	return parts[0], parts[1]
}

func LabelsToString(labels map[string]string) string {
	var labelsString []string
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := labels[key]
		if value == "" {
			labelsString = append(labelsString, key)
		} else {
			labelsString = append(labelsString, fmt.Sprintf("%s=%s", key, value))
		}
	}
	return strings.Join(labelsString, ", ")
}

func DescribeFormat(object interface{}, format string) error {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	t, err := template.New("").Parse(format)
	if err != nil {
		return err
	}
	return t.Execute(os.Stdout, object)
}

func DescribeJSON(object interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(object)
}

func LocationToSchema(location hcloud.Location) schema.Location {
	return schema.Location{
		ID:          location.ID,
		Name:        location.Name,
		Description: location.Description,
		Country:     location.Country,
		City:        location.City,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		NetworkZone: string(location.NetworkZone),
	}
}

func DatacenterToSchema(datacenter hcloud.Datacenter) schema.Datacenter {
	datacenterSchema := schema.Datacenter{
		ID:          datacenter.ID,
		Name:        datacenter.Name,
		Description: datacenter.Description,
		Location:    LocationToSchema(*datacenter.Location),
	}
	for _, st := range datacenter.ServerTypes.Supported {
		datacenterSchema.ServerTypes.Supported = append(datacenterSchema.ServerTypes.Supported, st.ID)
	}
	for _, st := range datacenter.ServerTypes.Available {
		datacenterSchema.ServerTypes.Available = append(datacenterSchema.ServerTypes.Available, st.ID)
	}
	return datacenterSchema
}

func ServerTypeToSchema(serverType hcloud.ServerType) schema.ServerType {
	serverTypeSchema := schema.ServerType{
		ID:          serverType.ID,
		Name:        serverType.Name,
		Description: serverType.Description,
		Cores:       serverType.Cores,
		Memory:      serverType.Memory,
		Disk:        serverType.Disk,
		StorageType: string(serverType.StorageType),
		CPUType:     string(serverType.CPUType),
	}
	for _, pricing := range serverType.Pricings {
		serverTypeSchema.Prices = append(serverTypeSchema.Prices, schema.PricingServerTypePrice{
			Location: pricing.Location.Name,
			PriceHourly: schema.Price{
				Net:   pricing.Hourly.Net,
				Gross: pricing.Hourly.Gross,
			},
			PriceMonthly: schema.Price{
				Net:   pricing.Monthly.Net,
				Gross: pricing.Monthly.Gross,
			},
		})
	}
	return serverTypeSchema
}

func ImageToSchema(image hcloud.Image) schema.Image {
	imageSchema := schema.Image{
		ID:          image.ID,
		Name:        hcloud.String(image.Name),
		Description: image.Description,
		Status:      string(image.Status),
		Type:        string(image.Type),
		ImageSize:   &image.ImageSize,
		DiskSize:    image.DiskSize,
		Created:     image.Created,
		OSFlavor:    image.OSFlavor,
		OSVersion:   hcloud.String(image.OSVersion),
		RapidDeploy: image.RapidDeploy,
		Protection: schema.ImageProtection{
			Delete: image.Protection.Delete,
		},
		Deprecated: image.Deprecated,
		Labels:     image.Labels,
	}
	if image.CreatedFrom != nil {
		imageSchema.CreatedFrom = &schema.ImageCreatedFrom{
			ID:   image.CreatedFrom.ID,
			Name: image.CreatedFrom.Name,
		}
	}
	if image.BoundTo != nil {
		imageSchema.BoundTo = hcloud.Int(image.BoundTo.ID)
	}
	return imageSchema
}

func ISOToSchema(iso hcloud.ISO) schema.ISO {
	return schema.ISO{
		ID:          iso.ID,
		Name:        iso.Name,
		Description: iso.Description,
		Deprecated:  iso.Deprecated,
	}
}

func LoadBalancerTypeToSchema(loadBalancerType hcloud.LoadBalancerType) schema.LoadBalancerType {
	loadBalancerTypeSchema := schema.LoadBalancerType{
		ID:                      loadBalancerType.ID,
		Name:                    loadBalancerType.Name,
		Description:             loadBalancerType.Description,
		MaxConnections:          loadBalancerType.MaxConnections,
		MaxServices:             loadBalancerType.MaxServices,
		MaxTargets:              loadBalancerType.MaxTargets,
		MaxAssignedCertificates: loadBalancerType.MaxAssignedCertificates,
	}
	for _, pricing := range loadBalancerType.Pricings {
		loadBalancerTypeSchema.Prices = append(loadBalancerTypeSchema.Prices, schema.PricingLoadBalancerTypePrice{
			Location: pricing.Location.Name,
			PriceHourly: schema.Price{
				Net:   pricing.Hourly.Net,
				Gross: pricing.Hourly.Gross,
			},
			PriceMonthly: schema.Price{
				Net:   pricing.Monthly.Net,
				Gross: pricing.Monthly.Gross,
			},
		})
	}
	return loadBalancerTypeSchema
}

func PlacementGroupToSchema(placementGroup hcloud.PlacementGroup) schema.PlacementGroup {
	return schema.PlacementGroup{
		ID:      placementGroup.ID,
		Name:    placementGroup.Name,
		Labels:  placementGroup.Labels,
		Created: placementGroup.Created,
		Type:    string(placementGroup.Type),
		Servers: placementGroup.Servers,
	}
}

// ValidateRequiredFlags ensures that flags has values for all flags with
// the passed names.
//
// This function duplicates the functionality cobra provides when calling
// MarkFlagRequired. However, in some cases a flag cannot be marked as required
// in cobra, for example when it depends on other flags. In those cases this
// function comes in handy.
func ValidateRequiredFlags(flags *pflag.FlagSet, names ...string) error {
	var missingFlags []string

	for _, name := range names {
		if !flags.Changed(name) {
			missingFlags = append(missingFlags, `"`+name+`"`)
		}
	}
	if len(missingFlags) > 0 {
		return fmt.Errorf("hcloud: required flag(s) %s not set", strings.Join(missingFlags, ", "))
	}
	return nil
}
