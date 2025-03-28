package primaryip_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.PrimaryIPClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.PrimaryIP{ID: 123}, nil, nil)
	fx.Client.PrimaryIPClient.EXPECT().
		Update(gomock.Any(), &hcloud.PrimaryIP{ID: 123}, hcloud.PrimaryIPUpdateOpts{
			Labels: hcloud.Ptr(map[string]string{
				"key": "value",
			}),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to Primary IP 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	primaryIP := &hcloud.PrimaryIP{
		ID: 123,
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.PrimaryIPClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(primaryIP, nil, nil)
	fx.Client.PrimaryIPClient.EXPECT().
		Update(gomock.Any(), primaryIP, hcloud.PrimaryIPUpdateOpts{
			Labels: hcloud.Ptr(make(map[string]string)),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from Primary IP 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
