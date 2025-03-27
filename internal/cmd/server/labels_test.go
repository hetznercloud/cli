package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Server{ID: 123}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), &hcloud.Server{ID: 123}, hcloud.ServerUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to Server 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestMultiLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.LabelCmds.AddCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Server{ID: 123}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), &hcloud.Server{ID: 123}, hcloud.ServerUpdateOpts{
			Labels: map[string]string{
				"foo": "bar",
				"baz": "qux",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "foo=bar", "baz=qux"})

	expOut := "Label(s) foo, baz added to Server 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{
		ID: 123,
		Labels: map[string]string{
			"key": "value",
		},
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), srv, hcloud.ServerUpdateOpts{
			Labels: make(map[string]string),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from Server 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestMultiLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{
		ID: 123,
		Labels: map[string]string{
			"key": "value",
			"foo": "bar",
			"baz": "qux",
		},
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), srv, hcloud.ServerUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "foo", "baz"})

	expOut := "Label(s) foo, baz removed from Server 123\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
