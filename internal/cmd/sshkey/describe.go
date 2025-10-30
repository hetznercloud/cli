package sshkey

import (
	"fmt"
	"io"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.SSHKey]{
	ResourceNameSingular: "SSH Key",
	ShortDescription:     "Describe an SSH Key",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.SSHKey().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.SSHKey, any, error) {
		key, _, err := s.Client().SSHKey().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return key, hcloud.SchemaFromSSHKey(key), nil
	},
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, sshKey *hcloud.SSHKey) error {
		_, _ = fmt.Fprintf(out, "ID:\t%d\n", sshKey.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", sshKey.Name)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(sshKey.Created), humanize.Time(sshKey.Created))
		_, _ = fmt.Fprintf(out, "Fingerprint:\t%s\n", sshKey.Fingerprint)
		_, _ = fmt.Fprintf(out, "Public Key:\n%s\n", strings.TrimSpace(sshKey.PublicKey))

		util.DescribeLabels(out, sshKey.Labels, "")

		return nil
	},
}
