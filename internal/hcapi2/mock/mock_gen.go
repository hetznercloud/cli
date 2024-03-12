package hcapi2_mock

//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_action_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ActionClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_certificate_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 CertificateClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_datacenter_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 DatacenterClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_image_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ImageClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_iso_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ISOClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_firewall_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 FirewallClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_floating_ip_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 FloatingIPClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_primary_ip_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 PrimaryIPClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_location_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 LocationClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_loadbalancer_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 LoadBalancerClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_loadbalancer_type_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 LoadBalancerTypeClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_network_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 NetworkClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_server_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ServerClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_server_type_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 ServerTypeClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_ssh_key_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 SSHKeyClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_volume_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 VolumeClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_placement_group_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 PlacementGroupClient
//go:generate go run github.com/golang/mock/mockgen -package hcapi2_mock -destination zz_rdns_client_mock.go github.com/hetznercloud/cli/internal/hcapi2 RDNSClient
