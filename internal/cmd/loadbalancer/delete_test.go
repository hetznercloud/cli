package loadbalancer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancer"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDelete(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	loadBalancer := &hcloud.LoadBalancer{
		ID:   123,
		Name: "test",
	}

	fx.Client.LoadBalancer.EXPECT().
		Get(gomock.Any(), "test").
		Return(loadBalancer, nil, nil)
	fx.Client.LoadBalancer.EXPECT().
		Delete(gomock.Any(), loadBalancer).
		Return(nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := "Load Balancer test deleted\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestDeleteMultiple(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancer.DeleteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	loadBalancers := []*hcloud.LoadBalancer{
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
	for _, lb := range loadBalancers {
		names = append(names, lb.Name)
		expOutBuilder.WriteString(fmt.Sprintf("Load Balancer %s deleted\n", lb.Name))
		fx.Client.LoadBalancer.EXPECT().
			Get(gomock.Any(), lb.Name).
			Return(lb, nil, nil)
		fx.Client.LoadBalancer.EXPECT().
			Delete(gomock.Any(), lb).
			Return(nil, nil)
	}

	out, errOut, err := fx.Run(cmd, names)
	expOut := expOutBuilder.String()

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}
