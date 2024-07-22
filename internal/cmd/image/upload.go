package image

import (
	"fmt"
	"github.com/apricote/hcloud-upload-image/hcloudimages"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

const (
	uploadFlagImageURL     = "image-url"
	uploadFlagImagePath    = "image-path"
	uploadFlagCompression  = "compression"
	uploadFlagArchitecture = "architecture"
	uploadFlagServerType   = "server-type"
	uploadFlagDescription  = "description"
	uploadFlagLabels       = "labels"
)

var UploadCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "upload (--image-path=<local-path> | --image-url=<url>) (--architecture=<x86|arm> | --server-type=<server-type>)",
			Short: "[alpha] Upload an OS image to a snapshot",
			Long: `This command implements a fake "upload", by going through a real server and snapshots.
This does cost a bit of money for the server.`,
		}

		cmd.Flags().String(uploadFlagImageURL, "", "Remote URL of the disk image that should be uploaded")
		cmd.Flags().String(uploadFlagImagePath, "", "Local path to the disk image that should be uploaded")
		cmd.MarkFlagsMutuallyExclusive(uploadFlagImageURL, uploadFlagImagePath)
		cmd.MarkFlagsOneRequired(uploadFlagImageURL, uploadFlagImagePath)

		cmd.Flags().String(uploadFlagCompression, "", "Type of compression that was used on the disk image [choices: bz2, xz]")
		_ = cmd.RegisterFlagCompletionFunc(
			uploadFlagCompression,
			cobra.FixedCompletions([]string{string(hcloudimages.CompressionBZ2), string(hcloudimages.CompressionXZ)}, cobra.ShellCompDirectiveNoFileComp),
		)

		cmd.Flags().String(uploadFlagArchitecture, "", "CPU architecture of the disk image [choices: x86, arm]")
		_ = cmd.RegisterFlagCompletionFunc(
			uploadFlagArchitecture,
			cobra.FixedCompletions([]string{string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)}, cobra.ShellCompDirectiveNoFileComp),
		)

		cmd.Flags().String(uploadFlagServerType, "", "Explicitly use this server type to generate the image. Mutually exclusive with --architecture.")

		// Only one of them needs to be set
		cmd.MarkFlagsOneRequired(uploadFlagArchitecture, uploadFlagServerType)
		cmd.MarkFlagsMutuallyExclusive(uploadFlagArchitecture, uploadFlagServerType)

		cmd.Flags().String(uploadFlagDescription, "", "Description for the resulting image")

		cmd.Flags().StringToString(uploadFlagLabels, map[string]string{}, "Labels for the resulting image")

		return cmd

	},
	Run: func(state state.State, cmd *cobra.Command, strings []string) error {
		hcloudClient := state.Client().(*hcapi2.ActualClient).Client
		imageClient := hcloudimages.NewClient(hcloudClient)

		imageURLString, _ := cmd.Flags().GetString(uploadFlagImageURL)
		imagePathString, _ := cmd.Flags().GetString(uploadFlagImagePath)
		imageCompression, _ := cmd.Flags().GetString(uploadFlagCompression)
		architecture, _ := cmd.Flags().GetString(uploadFlagArchitecture)
		serverType, _ := cmd.Flags().GetString(uploadFlagServerType)
		description, _ := cmd.Flags().GetString(uploadFlagDescription)
		labels, _ := cmd.Flags().GetStringToString(uploadFlagLabels)

		options := hcloudimages.UploadOptions{
			ImageCompression: hcloudimages.Compression(imageCompression),
			Description:      hcloud.Ptr(description),
			Labels:           labels,
		}

		if imageURLString != "" {
			imageURL, err := url.Parse(imageURLString)
			if err != nil {
				return fmt.Errorf("unable to parse url from --%s=%q: %w", uploadFlagImageURL, imageURLString, err)
			}

			options.ImageURL = imageURL
		} else if imagePathString != "" {
			imageFile, err := os.Open(imagePathString)
			if err != nil {
				return fmt.Errorf("unable to read file from --%s=%q: %w", uploadFlagImagePath, imagePathString, err)
			}

			options.ImageReader = imageFile
		}

		if architecture != "" {
			options.Architecture = hcloud.Architecture(architecture)
		} else if serverType != "" {
			options.ServerType = &hcloud.ServerType{Name: serverType}
		}

		image, err := imageClient.Upload(cmd.Context(), options)
		if err != nil {
			return fmt.Errorf("failed to upload the image: %w", err)
		}

		cmd.Printf("Uploaded as image %d", image.ID)
		return nil
	},
}
