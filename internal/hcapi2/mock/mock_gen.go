package mock

//go:generate mockgen -package mock -destination zz_action_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ActionClient
//go:generate mockgen -package mock -destination zz_certificate_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 CertificateClient
//go:generate mockgen -package mock -destination zz_datacenter_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 DatacenterClient
//go:generate mockgen -package mock -destination zz_image_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ImageClient
//go:generate mockgen -package mock -destination zz_iso_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ISOClient
//go:generate mockgen -package mock -destination zz_firewall_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 FirewallClient
//go:generate mockgen -package mock -destination zz_floating_ip_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 FloatingIPClient
//go:generate mockgen -package mock -destination zz_primary_ip_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 PrimaryIPClient
//go:generate mockgen -package mock -destination zz_location_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 LocationClient
//go:generate mockgen -package mock -destination zz_loadbalancer_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 LoadBalancerClient
//go:generate mockgen -package mock -destination zz_loadbalancer_type_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 LoadBalancerTypeClient
//go:generate mockgen -package mock -destination zz_network_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 NetworkClient
//go:generate mockgen -package mock -destination zz_server_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ServerClient
//go:generate mockgen -package mock -destination zz_server_type_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ServerTypeClient
//go:generate mockgen -package mock -destination zz_ssh_key_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 SSHKeyClient
//go:generate mockgen -package mock -destination zz_volume_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 VolumeClient
//go:generate mockgen -package mock -destination zz_placement_group_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 PlacementGroupClient
//go:generate mockgen -package mock -destination zz_rdns_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 RDNSClient
//go:generate mockgen -package mock -destination zz_pricing_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 PricingClient
//go:generate mockgen -package mock -destination zz_storage_box_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 StorageBoxClient
//go:generate mockgen -package mock -destination zz_storage_box_type_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 StorageBoxTypeClient
