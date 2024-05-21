package floatingip_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/floatingip"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	floatingIP := &hcloud.FloatingIP{
		ID:   123,
		Name: "test",
	}

	fx.Client.FloatingIP.EXPECT().
		Get(gomock.Any(), "test").
		Return(floatingIP, nil, nil)
	fx.Client.FloatingIP.EXPECT().
		Delete(gomock.Any(), floatingIP).
		Return(nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Floating IP test deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := floatingip.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	ips := []*hcloud.FloatingIP{
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
		expOutBuilder.WriteString(fmt.Sprintf("Floating IP %s deleted\n", ip.Name))
		fx.Client.FloatingIP.EXPECT().
			Get(gomock.Any(), ip.Name).
			Return(ip, nil, nil)
		fx.Client.FloatingIP.EXPECT().
			Delete(gomock.Any(), ip).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)
	expOut := expOutBuilder.String()

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
