package subaccount_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storagebox/subaccount"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	time.Local = time.UTC

	cmd := subaccount.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	sb := &hcloud.StorageBox{
		ID:   123,
		Name: "my-storage-box",
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

	fx.Client.StorageBoxClient.EXPECT().
		Get(gomock.Any(), "my-storage-box").
		Return(sb, nil, nil)
	fx.Client.StorageBoxClient.EXPECT().
		GetSubaccount(gomock.Any(), sb, "42").
		Return(sbs, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"my-storage-box", "42"})

	expOut := fmt.Sprintf(`ID:			42
Description:		host01 backup
Created:		Sat Jan 30 23:55:00 UTC 2016 (%s)
Username:		u1337-sub1
Home Directory:		my_backups/host01.my.company
Server:			u1337-sub1.your-storagebox.de
Access Settings:
  Reachable Externally:	false
  Samba Enabled:	false
  SSH Enabled:		false
  WebDAV Enabled:	false
  Readonly:		false
Labels:
  environment: prod
  example.com/my: label
  just-a-key: 
Storage Box:
  ID:			123
`, humanize.Time(sbs.Created))

	require.NoError(t, err)
	assert.Equal(t, ExperimentalWarning, errOut)
	assert.Equal(t, expOut, out)
}
