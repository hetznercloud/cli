package subaccount_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/subaccount"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestList(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := subaccount.ListCmd.CobraCommand(fx.State())

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "test",
	}
	sbs := &hcloud.StorageBoxSubaccount{
		ID:            42,
		Username:      "u1337-sub1",
		HomeDirectory: "my_backups/host01.my.company",
		Server:        "u1337-sub1.your-storagebox.de",
		AccessSettings: &hcloud.StorageBoxSubaccountAccessSettings{
			ReachableExternally: false,
			Readonly:            false,
			SambaEnabled:        false,
			SSHEnabled:          false,
			WebDAVEnabled:       false,
		},
		Description: "host01 backup",
		Labels: map[string]string{
			"environment":    "prod",
			"example.com/my": "label",
			"just-a-key":     "",
		},
		Created:    time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC),
		StorageBox: sb,
	}

	fx.ExpectEnsureToken()
	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "test").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		AllSubaccountsWithOpts(
			gomock.Any(),
			sb,
			hcloud.StorageBoxSubaccountListOpts{Sort: []string{"id:asc"}},
		).
		Return([]*hcloud.StorageBoxSubaccount{
			sbs,
		}, nil)

	out, errOut, err := fx.Run(cmd, []string{"test"})

	expOut := fmt.Sprintf(`ID   USERNAME     HOME DIRECTORY                 DESCRIPTION     SERVER                          AGE  
42   u1337-sub1   my_backups/host01.my.company   host01 backup   u1337-sub1.your-storagebox.de   %s
`, util.Age(sbs.Created, time.Now()))

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
