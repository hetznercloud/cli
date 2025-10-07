package zone_test

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

//go:embed testdata/example.zone
var zonefile string

func TestExportZonefile(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.ExportZonefileCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{ID: 123}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		ExportZonefile(gomock.Any(), z).
		Return(hcloud.ZoneExportZonefileResult{
			Zonefile: zonefile,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"123"})

	expOut := zonefile

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}

func TestExportZonefileJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.ExportZonefileCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{ID: 123}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		ExportZonefile(gomock.Any(), z).
		Return(hcloud.ZoneExportZonefileResult{
			Zonefile: zonefile,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "-o=json"})

	expOut := fmt.Sprintf(`{ "zonefile": %q }`, zonefile)

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.JSONEq(t, expOut, out)
}

func TestExportZonefileYAML(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.ExportZonefileCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{ID: 123}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		ExportZonefile(gomock.Any(), z).
		Return(hcloud.ZoneExportZonefileResult{
			Zonefile: zonefile,
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "-o=yaml"})

	// This will work since JSON is a subset of YAML
	expOut := fmt.Sprintf(`{ "zonefile": %q }`, zonefile)

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.YAMLEq(t, expOut, out)
}
