package server_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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

	out, _, err := fx.Run(cmd, []string{"123", "key=value"})

	expOut := "Label(s) key added to server 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestMultiLabelAdd(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.AddCobraCommand(fx.State())
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

	out, _, err := fx.Run(cmd, []string{"123", "foo=bar", "baz=qux"})

	expOut := "Label(s) foo, baz added to server 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Server{
			ID: 123,
			Labels: map[string]string{
				"key": "value",
			},
		}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), &hcloud.Server{ID: 123}, hcloud.ServerUpdateOpts{
			Labels: make(map[string]string),
		})

	out, _, err := fx.Run(cmd, []string{"123", "key"})

	expOut := "Label(s) key removed from server 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestMultiLabelRemove(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := LabelCmds.RemoveCobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Server{
			ID: 123,
			Labels: map[string]string{
				"key": "value",
				"foo": "bar",
				"baz": "qux",
			},
		}, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Update(gomock.Any(), &hcloud.Server{ID: 123}, hcloud.ServerUpdateOpts{
			Labels: map[string]string{
				"key": "value",
			},
		})

	out, _, err := fx.Run(cmd, []string{"123", "foo", "baz"})

	expOut := "Label(s) foo, baz removed from server 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
