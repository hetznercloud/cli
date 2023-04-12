package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
)

// CreateCmd defines a command for creating a server.
var CreateCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "create FLAGS",
			Short: "Create a server",
		}

		cmd.Flags().String("name", "", "Server name (required)")
		cmd.MarkFlagRequired("name")

		cmd.Flags().String("type", "", "Server type (ID or name) (required)")
		cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidatesF(client.ServerType().Names))
		cmd.MarkFlagRequired("type")

		cmd.Flags().String("image", "", "Image (ID or name) (required)")
		cmd.RegisterFlagCompletionFunc("image", cmpl.SuggestCandidatesF(client.Image().Names))
		cmd.MarkFlagRequired("image")

		cmd.Flags().String("location", "", "Location (ID or name)")
		cmd.RegisterFlagCompletionFunc("location", cmpl.SuggestCandidatesF(client.Location().Names))

		cmd.Flags().String("datacenter", "", "Datacenter (ID or name)")
		cmd.RegisterFlagCompletionFunc("datacenter", cmpl.SuggestCandidatesF(client.Datacenter().Names))

		cmd.Flags().StringSlice("ssh-key", nil, "ID or name of SSH key to inject (can be specified multiple times)")
		cmd.RegisterFlagCompletionFunc("ssh-key", cmpl.SuggestCandidatesF(client.SSHKey().Names))

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		cmd.Flags().StringArray("user-data-from-file", []string{}, "Read user data from specified file (use - to read from stdin)")

		cmd.Flags().Bool("start-after-create", true, "Start server right after creation")

		cmd.Flags().StringSlice("volume", nil, "ID or name of volume to attach (can be specified multiple times)")
		cmd.RegisterFlagCompletionFunc("volume", cmpl.SuggestCandidatesF(client.Volume().Names))

		cmd.Flags().StringSlice("network", nil, "ID or name of network to attach the server to (can be specified multiple times)")
		cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(client.Network().Names))

		cmd.Flags().StringSlice("firewall", nil, "ID or name of Firewall to attach the server to (can be specified multiple times)")
		cmd.RegisterFlagCompletionFunc("firewall", cmpl.SuggestCandidatesF(client.Firewall().Names))

		cmd.Flags().Bool("automount", false, "Automount volumes after attach (default: false)")
		cmd.Flags().Bool("allow-deprecated-image", false, "Enable the use of deprecated images (default: false)")

		cmd.Flags().String("placement-group", "", "Placement Group (ID of name)")
		cmd.RegisterFlagCompletionFunc("placement-group", cmpl.SuggestCandidatesF(client.PlacementGroup().Names))
		cmd.Flags().String("primary-ipv4", "", "Primary IPv4 (ID of name)")
		cmd.RegisterFlagCompletionFunc("primary-ipv4", cmpl.SuggestCandidatesF(client.PrimaryIP().IPv4Names))
		cmd.Flags().String("primary-ipv6", "", "Primary IPv6 (ID of name)")
		cmd.RegisterFlagCompletionFunc("primary-ipv6", cmpl.SuggestCandidatesF(client.PrimaryIP().IPv6Names))

		cmd.Flags().Bool("without-ipv4", false, "Creates the server without an IPv4 (default: false)")
		cmd.Flags().Bool("without-ipv6", false, "Creates the server without an IPv6 (default: false)")
		return cmd
	},

	Run: func(ctx context.Context, client hcapi2.Client, actionWaiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		opts, err := createOptsFromFlags(ctx, client, cmd.Flags())
		if err != nil {
			return err
		}

		result, _, err := client.Server().Create(ctx, opts)
		if err != nil {
			return err
		}

		if err := actionWaiter.ActionProgress(ctx, result.Action); err != nil {
			return err
		}
		if err := actionWaiter.WaitForActions(ctx, result.NextActions); err != nil {
			return err
		}

		server, _, err := client.Server().GetByID(ctx, result.Server.ID)
		if err != nil {
			return err
		}
		fmt.Printf("Server %d created\n", result.Server.ID)
		if !server.PublicNet.IPv4.IsUnspecified() {
			fmt.Printf("IPv4: %s\n", server.PublicNet.IPv4.IP.String())
		}
		if !server.PublicNet.IPv6.IsUnspecified() {
			fmt.Printf("IPv6: %s1\n", server.PublicNet.IPv6.Network.IP.String())
			fmt.Printf("IPv6 Network: %s\n", server.PublicNet.IPv6.Network.String())
		}
		if len(server.PrivateNet) > 0 {
			var networks []string
			for _, network := range server.PrivateNet {
				networks = append(networks, fmt.Sprintf("- %s (%s)", network.IP.String(), client.Network().Name(network.Network.ID)))
			}
			fmt.Printf("Private Networks:\n\t%s\n", strings.Join(networks, "\n"))
		}
		// Only print the root password if it's not empty,
		// which is only the case if it wasn't created with an SSH key.
		if result.RootPassword != "" {
			fmt.Printf("Root password: %s\n", result.RootPassword)
		}

		return nil
	},
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

func createOptsFromFlags(
	ctx context.Context, client hcapi2.Client, flags *pflag.FlagSet,
) (opts hcloud.ServerCreateOpts, err error) {
	name, _ := flags.GetString("name")
	serverTypeName, _ := flags.GetString("type")
	imageIDorName, _ := flags.GetString("image")
	location, _ := flags.GetString("location")
	datacenter, _ := flags.GetString("datacenter")
	userDataFiles, _ := flags.GetStringArray("user-data-from-file")
	startAfterCreate, _ := flags.GetBool("start-after-create")
	sshKeys, _ := flags.GetStringSlice("ssh-key")
	labels, _ := flags.GetStringToString("label")
	volumes, _ := flags.GetStringSlice("volume")
	networks, _ := flags.GetStringSlice("network")
	firewalls, _ := flags.GetStringSlice("firewall")
	automount, _ := flags.GetBool("automount")
	allowDeprecatedImage, _ := flags.GetBool("allow-deprecated-image")
	placementGroupIDorName, _ := flags.GetString("placement-group")
	withoutIPv4, _ := flags.GetBool("without-ipv4")
	withoutIPv6, _ := flags.GetBool("without-ipv6")
	primaryIPv4IDorName, _ := flags.GetString("primary-ipv4")
	primaryIPv6IDorName, _ := flags.GetString("primary-ipv6")

	serverType, _, err := client.ServerType().Get(ctx, serverTypeName)
	if err != nil {
		return
	}

	// Select correct image based on server type architecture
	image, _, err := client.Image().GetForArchitecture(ctx, imageIDorName, serverType.Architecture)
	if err != nil {
		return
	}

	if image == nil {
		err = fmt.Errorf("image %s for architecture %s not found", imageIDorName, serverType.Architecture)
		return
	}

	if !image.Deprecated.IsZero() {
		if allowDeprecatedImage {
			fmt.Printf("Attention: image %s is deprecated. It will continue to be available until %s.\n", image.Name, image.Deprecated.AddDate(0, 3, 0).Format("2006-01-02"))
		} else {
			err = fmt.Errorf("image %s is deprecated, please use --allow-deprecated-image to create a server with this image. It will continue to be available until %s", image.Name, image.Deprecated.AddDate(0, 3, 0).Format("2006-01-02"))
			return
		}
	}

	if withoutIPv4 && withoutIPv6 && len(networks) == 0 {
		err = fmt.Errorf("a server can not be created without IPv4, IPv6 and a private network. Choose at least one of those options to create the server")
		return
	}
	opts = hcloud.ServerCreateOpts{
		Name:             name,
		ServerType:       serverType,
		Image:            image,
		Labels:           labels,
		StartAfterCreate: &startAfterCreate,
		Automount:        &automount,
	}
	publicNetConfiguration := &hcloud.ServerCreatePublicNet{}
	if !withoutIPv4 {
		publicNetConfiguration.EnableIPv4 = true
	}
	if !withoutIPv6 {
		publicNetConfiguration.EnableIPv6 = true
	}
	if primaryIPv4IDorName != "" {
		var primaryIPv4 *hcloud.PrimaryIP
		primaryIPv4, _, err = client.PrimaryIP().Get(ctx, primaryIPv4IDorName)
		if err != nil {
			return
		}
		if primaryIPv4 == nil {
			err = fmt.Errorf("primary ipv4 not found: %s", primaryIPv4IDorName)
			return
		}
		publicNetConfiguration.IPv4 = primaryIPv4
	}
	if primaryIPv6IDorName != "" {
		var primaryIPv6 *hcloud.PrimaryIP
		primaryIPv6, _, err = client.PrimaryIP().Get(ctx, primaryIPv6IDorName)
		if err != nil {
			return
		}
		if primaryIPv6 == nil {
			err = fmt.Errorf("primary ipv4 not found: %s", primaryIPv6IDorName)
			return
		}
		publicNetConfiguration.IPv6 = primaryIPv6
	}
	opts.PublicNet = publicNetConfiguration
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
		sshKey, _, err = client.SSHKey().Get(ctx, sshKeyIDOrName)
		if err != nil {
			return
		}

		if sshKey == nil {
			sshKey, err = getSSHKeyForFingerprint(ctx, client, sshKeyIDOrName)
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
		volume, _, err = client.Volume().Get(ctx, volumeIDOrName)
		if err != nil {
			return
		}

		if volume == nil {
			err = fmt.Errorf("volume not found: %s", volumeIDOrName)
			return
		}
		opts.Volumes = append(opts.Volumes, volume)
	}
	for _, networkIDOrName := range networks {
		var network *hcloud.Network
		network, _, err = client.Network().Get(ctx, networkIDOrName)
		if err != nil {
			return
		}

		if network == nil {
			err = fmt.Errorf("network not found: %s", networkIDOrName)
			return
		}
		opts.Networks = append(opts.Networks, network)
	}
	for _, firewallIDOrName := range firewalls {
		var firewall *hcloud.Firewall
		firewall, _, err = client.Firewall().Get(ctx, firewallIDOrName)
		if err != nil {
			return
		}

		if firewall == nil {
			err = fmt.Errorf("firewall not found: %s", firewallIDOrName)
			return
		}
		opts.Firewalls = append(opts.Firewalls, &hcloud.ServerCreateFirewall{Firewall: *firewall})
	}

	if datacenter != "" {
		opts.Datacenter = &hcloud.Datacenter{Name: datacenter}
	}
	if location != "" {
		opts.Location = &hcloud.Location{Name: location}
	}
	if placementGroupIDorName != "" {
		var placementGroup *hcloud.PlacementGroup
		placementGroup, _, err = client.PlacementGroup().Get(ctx, placementGroupIDorName)
		if err != nil {
			return
		}
		if placementGroup == nil {
			err = fmt.Errorf("placement group not found: %s", placementGroupIDorName)
			return
		}
		opts.PlacementGroup = placementGroup
	}

	return
}

func getSSHKeyForFingerprint(
	ctx context.Context, client hcapi2.Client, file string,
) (sshKey *hcloud.SSHKey, err error) {
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
	sshKey, _, err = client.SSHKey().GetByFingerprint(ctx, ssh.FingerprintLegacyMD5(publicKey))
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
