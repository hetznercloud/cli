package primaryip_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/primaryip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.DeleteCmd.CobraCommand(fx.State())
	ip := &hcloud.PrimaryIP{ID: 13}
	fx.ExpectEnsureToken()
	fx.Client.PrimaryIPClient.EXPECT().
		Get(
			gomock.Any(),
			"13",
		).
		Return(
			ip,
			&hcloud.Response{},
			nil,
		)
	fx.Client.PrimaryIPClient.EXPECT().
		Delete(
			gomock.Any(),
			ip,
		).
		Return(
			&hcloud.Response{},
			nil,
		)

	out, errOut, err := fx.Run(cmd, []string{"13"})

	expOut := "Primary IP 13 deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := primaryip.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	ips := []*hcloud.PrimaryIP{
		{
			ID:   123,
			Name: "test1",
		},
		{
			ID:   456,
			Name: "test2",
		},
		{
			ID:   789,
			Name: "test3",
		},
	}

	expOutBuilder := strings.Builder{}

	var names []string
	for _, ip := range ips {
		names = append(names, ip.Name)
		expOutBuilder.WriteString(fmt.Sprintf("Primary IP %s deleted\n", ip.Name))
		fx.Client.PrimaryIPClient.EXPECT().
			Get(gomock.Any(), ip.Name).
			Return(ip, nil, nil)
		fx.Client.PrimaryIPClient.EXPECT().
			Delete(gomock.Any(), ip).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)
	expOut := expOutBuilder.String()

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
