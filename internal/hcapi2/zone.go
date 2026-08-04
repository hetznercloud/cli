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
	Names(context.Context) ([]string, error)
	LabelKeys(context.Context, string) ([]string, error)
	RRSetLabelKeys(context.Context, string, string, hcloud.ZoneRRSetType) ([]string, error)
}

func NewZoneClient(client *hcloud.ZoneClient) ZoneClient {
	return &zoneClient{
		IZoneClient: client,
	}
}

type zoneClient struct {
	hcloud.IZoneClient
}

func (c *zoneClient) Names(ctx context.Context) ([]string, error) {
	zones, err := c.All(ctx)
	if err != nil {
		return nil, err
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
	return names, nil
}

func (c *zoneClient) LabelKeys(ctx context.Context, idOrName string) ([]string, error) {
	idOrName, err := util.ParseZoneIDOrName(idOrName)
	if err != nil {
		return nil, err
	}

	zone, _, err := c.Get(ctx, idOrName)
	if err != nil {
		return nil, err
	}
	if zone == nil {
		return nil, nil
	}
	return labelKeys(zone.Labels), nil
}

func (c *zoneClient) RRSetLabelKeys(ctx context.Context, zoneIDOrName, rrsetName string, rrsetType hcloud.ZoneRRSetType) ([]string, error) {
	zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
	if err != nil {
		return nil, err
	}

	rrset, _, err := c.GetRRSetByNameAndType(ctx, &hcloud.Zone{Name: zoneIDOrName}, rrsetName, rrsetType)
	if err != nil {
		return nil, err
	}
	if rrset == nil {
		return nil, nil
	}
	return labelKeys(rrset.Labels), nil
}
