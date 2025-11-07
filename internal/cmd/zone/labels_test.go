package zone_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/zone"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		ID:   123,
		Name: "example.com",
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		Update(gomock.Any(), z, hcloud.ZoneUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to Zone example.com\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := zone.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	z := &hcloud.Zone{
		ID:   123,
		Name: "example.com",
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.ZoneClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(z, nil, nil)
	fx.Client.ZoneClient.EXPECT().
		Update(gomock.Any(), z, hcloud.ZoneUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from Zone example.com\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
