package hcapi2

//go:generate mockgen -package hcapi2 -destination zz_datacenter_client_mock.go . DatacenterClient
//go:generate mockgen -package hcapi2 -destination zz_image_client_mock.go . ImageClient
//go:generate mockgen -package hcapi2 -destination zz_firewall_client_mock.go . FirewallClient
//go:generate mockgen -package hcapi2 -destination zz_location_client_mock.go . LocationClient
//go:generate mockgen -package hcapi2 -destination zz_network_client_mock.go . NetworkClient
//go:generate mockgen -package hcapi2 -destination zz_server_client_mock.go . ServerClient
//go:generate mockgen -package hcapi2 -destination zz_server_type_client_mock.go . ServerTypeClient
//go:generate mockgen -package hcapi2 -destination zz_ssh_key_client_mock.go . SSHKeyClient
//go:generate mockgen -package hcapi2 -destination zz_volume_client_mock.go . VolumeClient
