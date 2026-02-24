#!/usr/bin/env bash
mockgen -package mock -destination zz_hcapi_mock.go \
  github.com/hetznercloud/cli/internal/hcapi2 \
ActionClient,\
CertificateClient,\
DatacenterClient,\
ImageClient,\
ISOClient,\
FirewallClient,\
FloatingIPClient,\
PrimaryIPClient,\
LocationClient,\
LoadBalancerClient,\
LoadBalancerTypeClient,\
NetworkClient,\
ServerClient,\
ServerTypeClient,\
SSHKeyClient,\
VolumeClient,\
PlacementGroupClient,\
RDNSClient,\
PricingClient,\
StorageBoxClient,\
StorageBoxTypeClient,\
ZoneClient
