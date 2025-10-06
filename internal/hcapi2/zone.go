package hcapi2

import (
	"context"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// ZoneClient embeds the Hetzner Cloud Zone client and provides some
// additional helper functions.
type ZoneClient interface {
	hcloud.IZoneClient
	Names() []string
	LabelKeys(idOrName string) []string
	RRSetLabelKeys(zoneIDOrName, rrsetName string, rrsetType hcloud.ZoneRRSetType) []string
}

func NewZoneClient(client *hcloud.ZoneClient) ZoneClient {
	return &zoneClient{
		IZoneClient: client,
	}
}

type zoneClient struct {
	hcloud.IZoneClient
}

func (c *zoneClient) Names() []string {
	zones, err := c.All(context.Background())
	if err != nil || len(zones) == 0 {
		return nil
	}
	names := make([]string, 0, len(zones)*2)
	for _, zone := range zones {
		if zone.Name == "" {
			names = append(names, strconv.FormatInt(zone.ID, 10))
			continue
		}

		// Name as the API specifies it (IDNA-encoded)
		names = append(names, zone.Name)

		displayName := util.DisplayZoneName(zone.Name)
		if zone.Name != displayName {
			names = append(names, displayName)
		}
	}
	return names
}

func (c *zoneClient) LabelKeys(idOrName string) []string {
	if idOrIDNAName, err := util.ParseZoneIDOrName(idOrName); err == nil {
		// Ignore any errors in conversion
		idOrName = idOrIDNAName
	}

	zone, _, err := c.Get(context.Background(), idOrName)
	if err != nil || len(zone.Labels) == 0 {
		return nil
	}
	return labelKeys(zone.Labels)
}

func (c *zoneClient) RRSetLabelKeys(zoneIDOrName, rrsetName string, rrsetType hcloud.ZoneRRSetType) []string {
	if idOrIDNAName, err := util.ParseZoneIDOrName(zoneIDOrName); err == nil {
		// Ignore any errors in conversion
		zoneIDOrName = idOrIDNAName
	}

	rrset, _, err := c.GetRRSetByNameAndType(context.Background(), &hcloud.Zone{Name: zoneIDOrName}, rrsetName, rrsetType)
	if err != nil || len(rrset.Labels) == 0 {
		return nil
	}
	return labelKeys(rrset.Labels)
}
