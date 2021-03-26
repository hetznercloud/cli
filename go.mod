module github.com/hetznercloud/cli

require (
	github.com/cheggaaa/pb/v3 v3.0.7
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/structs v1.1.0
	github.com/guptarohit/asciigraph v0.5.1
<<<<<<< HEAD
	github.com/hetznercloud/hcloud-go v1.24.0
	github.com/pelletier/go-toml v1.8.0
	github.com/spf13/cobra v1.1.1
=======
	github.com/hetznercloud/hcloud-go v1.23.1
	github.com/pelletier/go-toml v1.8.1
	github.com/spf13/cobra v1.1.3
>>>>>>> 88a60a2 (Update dependencies)
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)

replace github.com/hetznercloud/hcloud-go => hetzner.cloud/integrations/hcloud-go v1.25.0-rc.3 // TODO: Remove before release

go 1.16
