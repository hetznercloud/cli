package hcapi2

//go:generate mockgen -package hcapi2 -destination zz_certificate_client_mock.go . CertificateClient
//go:generate mockgen -package hcapi2 -destination zz_datacenter_client_mock.go . DatacenterClient
//go:generate mockgen -package hcapi2 -destination zz_image_client_mock.go . ImageClient
//go:generate mockgen -package hcapi2 -destination zz_iso_client_mock.go . ISOClient
//go:generate mockgen -package hcapi2 -destination zz_firewall_client_mock.go . FirewallClient
//go:generate mockgen -package hcapi2 -destination zz_floating_ip_client_mock.go . FloatingIPClient
//go:generate mockgen -package hcapi2 -destination zz_primary_ip_client_mock.go . PrimaryIPClient
//go:generate mockgen -package hcapi2 -destination zz_location_client_mock.go . LocationClient
//go:generate mockgen -package hcapi2 -destination zz_loadbalancer_client_mock.go . LoadBalancerClient
//go:generate mockgen -package hcapi2 -destination zz_loadbalancer_type_client_mock.go . LoadBalancerTypeClient
//go:generate mockgen -package hcapi2 -destination zz_network_client_mock.go . NetworkClient
//go:generate mockgen -package hcapi2 -destination zz_server_client_mock.go . ServerClient
//go:generate mockgen -package hcapi2 -destination zz_server_type_client_mock.go . ServerTypeClient
//go:generate mockgen -package hcapi2 -destination zz_ssh_key_client_mock.go . SSHKeyClient
//go:generate mockgen -package hcapi2 -destination zz_volume_client_mock.go . VolumeClient
//go:generate mockgen -package hcapi2 -destination zz_placement_group_client_mock.go . PlacementGroupClient
//go:generate mockgen -package hcapi2 -destination zz_rdns_client_mock.go . RDNSClient
