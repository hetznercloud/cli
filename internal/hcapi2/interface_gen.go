package hcapi2

//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.DatacenterClient -as hcapi2.DatacenterClientBase -o zz_datacenter_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.ImageClient -as hcapi2.ImageClientBase -o zz_image_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.FirewallClient -as hcapi2.FirewallClientBase -o zz_firewall_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.LocationClient -as hcapi2.LocationClientBase -o zz_location_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.NetworkClient -as hcapi2.NetworkClientBase -o zz_network_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.ServerClient -as hcapi2.ServerClientBase -o zz_server_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.ServerTypeClient -as hcapi2.ServerTypeClientBase -o zz_server_type_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.SSHKeyClient -as hcapi2.SSHKeyClientBase -o zz_ssh_key_client_base.go
//go:generate interfacer -for github.com/hetznercloud/hcloud-go/hcloud.VolumeClient -as hcapi2.VolumeClientBase -o zz_volume_client_base.go
