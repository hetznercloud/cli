package cli

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
)

func newServerCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create FLAGS",
		Short:                 "Create a server",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerCreate),
	}
	cmd.Flags().String("name", "", "Server name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().String("type", "", "Server type (ID or name)")
	cmd.Flag("type").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_servertype_names"},
	}
	cmd.MarkFlagRequired("type")

	cmd.Flags().String("image", "", "Image (ID or name)")
	cmd.Flag("image").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_image_names"},
	}
	cmd.MarkFlagRequired("image")

	cmd.Flags().String("location", "", "Location (ID or name)")
	cmd.Flag("location").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_location_names"},
	}

	cmd.Flags().String("datacenter", "", "Datacenter (ID or name)")
	cmd.Flag("datacenter").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_datacenter_names"},
	}

	cmd.Flags().StringSlice("ssh-key", nil, "ID or name of SSH key to inject (can be specified multiple times)")
	cmd.Flag("ssh-key").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_sshkey_names"},
	}

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	cmd.Flags().StringArray("user-data-from-file", []string{}, "Read user data from specified file (use - to read from stdin)")

	cmd.Flags().Bool("start-after-create", true, "Start server right after creation (default: true)")

	cmd.Flags().StringSlice("volume", nil, "ID or name of volume to attach (can be specified multiple times)")
	cmd.Flag("volume").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_volume_names"},
	}
	cmd.Flags().StringSlice("network", nil, "ID of network to attach the server to (can be specified multiple times)")
	cmd.Flags().Bool("automount", false, "Automount volumes after attach (default: false)")
	return cmd
}

func runServerCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	opts, err := optsFromFlags(cli, cmd.Flags())
	if err != nil {
		return err
	}

	result, _, err := cli.Client().Server.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}
	if err := cli.WaitForActions(cli.Context, result.NextActions); err != nil {
		return err
	}

	fmt.Printf("Server %d created\n", result.Server.ID)
	fmt.Printf("IPv4: %s\n", result.Server.PublicNet.IPv4.IP.String())

	// Only print the root password if it's not empty,
	// which is only the case if it wasn't created with an SSH key.
	if result.RootPassword != "" {
		fmt.Printf("Root password: %s\n", result.RootPassword)
	}

	return nil
}

var userDataContentTypes = map[string]string{
	"#!":              "text/x-shellscript",
	"#include":        "text/x-include-url",
	"#cloud-config":   "text/cloud-config",
	"#upstart-job":    "text/upstart-job",
	"#cloud-boothook": "text/cloud-boothook",
	"#part-handler":   "text/part-handler",
}

func detectContentType(data string) string {
	for prefix, contentType := range userDataContentTypes {
		if strings.HasPrefix(data, prefix) {
			return contentType
		}
	}
	return ""
}

func buildUserData(files []string) (string, error) {
	var (
		buf = new(bytes.Buffer)
		mp  = multipart.NewWriter(buf)
	)

	fmt.Fprint(buf, "MIME-Version: 1.0\r\n")
	fmt.Fprint(buf, "Content-Type: multipart/mixed; boundary="+mp.Boundary()+"\r\n\r\n")

	for _, file := range files {
		var (
			data []byte
			err  error
		)
		if file == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(file)
		}
		if err != nil {
			return "", err
		}

		contentType := detectContentType(string(data))
		if contentType == "" {
			return "", fmt.Errorf("cannot detect user data type of file %q", file)
		}

		header := textproto.MIMEHeader{}
		header.Set("Content-Type", contentType)
		header.Set("Content-Transfer-Encoding", "base64")

		w, err := mp.CreatePart(header)
		if err != nil {
			return "", fmt.Errorf("failed to create multipart for file %q: %s", file, err)
		}

		if _, err := base64.NewEncoder(base64.StdEncoding, w).Write(data); err != nil {
			return "", fmt.Errorf("failed to encode data for file %q: %s", file, err)
		}
	}

	if err := mp.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func optsFromFlags(cli *CLI, flags *pflag.FlagSet) (opts hcloud.ServerCreateOpts, err error) {
	name, _ := flags.GetString("name")
	serverType, _ := flags.GetString("type")
	image, _ := flags.GetString("image")
	location, _ := flags.GetString("location")
	datacenter, _ := flags.GetString("datacenter")
	userDataFiles, _ := flags.GetStringArray("user-data-from-file")
	startAfterCreate, _ := flags.GetBool("start-after-create")
	sshKeys, _ := flags.GetStringSlice("ssh-key")
	labels, _ := flags.GetStringToString("label")
	volumes, _ := flags.GetStringSlice("volume")
	networks, _ := flags.GetStringSlice("network")
	automount, _ := flags.GetBool("automount")

	opts = hcloud.ServerCreateOpts{
		Name: name,
		ServerType: &hcloud.ServerType{
			Name: serverType,
		},
		Image: &hcloud.Image{
			Name: image,
		},
		Labels:           labels,
		StartAfterCreate: &startAfterCreate,
		Automount:        &automount,
	}

	if len(userDataFiles) == 1 {
		var data []byte
		if userDataFiles[0] == "-" {
			data, err = ioutil.ReadAll(os.Stdin)
		} else {
			data, err = ioutil.ReadFile(userDataFiles[0])
		}
		if err != nil {
			return
		}
		opts.UserData = string(data)
	} else if len(userDataFiles) > 1 {
		opts.UserData, err = buildUserData(userDataFiles)
		if err != nil {
			return
		}
	}

	for _, sshKeyIDOrName := range sshKeys {
		var sshKey *hcloud.SSHKey
		sshKey, _, err = cli.Client().SSHKey.Get(cli.Context, sshKeyIDOrName)
		if err != nil {
			return
		}

		if sshKey == nil {
			sshKey, err = getSSHKeyForFingerprint(cli, sshKeyIDOrName)
			if err != nil {
				return
			}
		}

		if sshKey == nil {
			err = fmt.Errorf("SSH key not found: %s", sshKeyIDOrName)
			return
		}
		opts.SSHKeys = append(opts.SSHKeys, sshKey)
	}
	for _, volumeIDOrName := range volumes {
		var volume *hcloud.Volume
		volume, _, err = cli.Client().Volume.Get(cli.Context, volumeIDOrName)
		if err != nil {
			return
		}

		if volume == nil {
			err = fmt.Errorf("volume not found: %s", volumeIDOrName)
			return
		}
		opts.Volumes = append(opts.Volumes, volume)
	}
	for _, networkID := range networks {
		var network *hcloud.Network
		network, _, err = cli.Client().Network.Get(cli.Context, networkID)
		if err != nil {
			return
		}

		if network == nil {
			err = fmt.Errorf("network not found: %s", networkID)
			return
		}
		opts.Networks = append(opts.Networks, network)
	}

	if datacenter != "" {
		opts.Datacenter = &hcloud.Datacenter{Name: datacenter}
	}
	if location != "" {
		opts.Location = &hcloud.Location{Name: location}
	}

	return
}

func getSSHKeyForFingerprint(cli *CLI, file string) (sshKey *hcloud.SSHKey, err error) {
	var (
		fileContent []byte
		publicKey   ssh.PublicKey
	)

	if fileContent, err = ioutil.ReadFile(file); err == os.ErrNotExist {
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("lookup SSH key by fingerprint: %v", err)
		return
	}

	if publicKey, _, _, _, err = ssh.ParseAuthorizedKey(fileContent); err != nil {
		err = fmt.Errorf("lookup SSH key by fingerprint: %v", err)
		return
	}
	sshKey, _, err = cli.Client().SSHKey.GetByFingerprint(cli.Context, ssh.FingerprintLegacyMD5(publicKey))
	if err != nil {
		err = fmt.Errorf("lookup SSH key by fingerprint: %v", err)
		return
	}
	if sshKey == nil {
		err = fmt.Errorf("SSH key not found by using fingerprint of public key: %s", file)
		return
	}
	return
}
