package hcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestActionFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"command": "create_server",
		"status": "success",
		"progress": 100,
		"started": "2016-01-30T23:55:00Z",
		"finished": "2016-01-30T23:56:13Z",
		"resources": [
			{
				"id": 42,
				"type": "server"
			}
		],
		"error": {
			"code": "action_failed",
			"message": "Action failed"
		}
	}`)

	var s schema.Action
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	action := ActionFromSchema(s)

	if action.ID != 1 {
		t.Errorf("unexpected ID: %v", action.ID)
	}
	if action.Command != "create_server" {
		t.Errorf("unexpected command: %v", action.Command)
	}
	if action.Status != "success" {
		t.Errorf("unexpected status: %v", action.Status)
	}
	if action.Progress != 100 {
		t.Errorf("unexpected progress: %d", action.Progress)
	}
	if !action.Started.Equal(time.Date(2016, 1, 30, 23, 55, 0, 0, time.UTC)) {
		t.Errorf("unexpected started: %v", action.Started)
	}
	if !action.Finished.Equal(time.Date(2016, 1, 30, 23, 56, 13, 0, time.UTC)) {
		t.Errorf("unexpected finished: %v", action.Started)
	}
	if action.ErrorCode != "action_failed" {
		t.Errorf("unexpected error code: %v", action.ErrorCode)
	}
	if action.ErrorMessage != "Action failed" {
		t.Errorf("unexpected error message: %v", action.ErrorMessage)
	}
}

func TestFloatingIPFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"description": "Web Frontend",
		"ip": "131.232.99.1",
		"type": "ipv4",
		"server": 42,
		"dns_ptr": "fip01.example.com",
		"home_location": {
			"id": 1,
			"name": "fsn1",
			"description": "Falkenstein DC Park 1",
			"country": "DE",
			"city": "Falkenstein",
			"latitude": 50.47612,
			"longitude": 12.370071
		}
	}`)

	var s schema.FloatingIP
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	floatingIP := FloatingIPFromSchema(s)

	if floatingIP.ID != 4711 {
		t.Errorf("unexpected ID: %v", floatingIP.ID)
	}
	if floatingIP.Description != "Web Frontend" {
		t.Errorf("unexpected description: %v", floatingIP.Description)
	}
	if floatingIP.IP != "131.232.99.1" {
		t.Errorf("unexpected IP: %v", floatingIP.IP)
	}
	if floatingIP.Type != FloatingIPTypeIPv4 {
		t.Errorf("unexpected type: %v", floatingIP.Type)
	}
	if floatingIP.Server == nil || floatingIP.Server.ID != 42 {
		t.Errorf("unexpected server: %v", floatingIP.Server)
	}
	if floatingIP.DNSPtr == nil || floatingIP.DNSPtr[floatingIP.IP] != "fip01.example.com" {
		t.Errorf("unexpected DNS ptr: %v", floatingIP.DNSPtr)
	}
	if floatingIP.HomeLocation == nil || floatingIP.HomeLocation.ID != 1 {
		t.Errorf("unexpected home location: %v", floatingIP.HomeLocation)
	}
}

func TestISOFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
		"description": "FreeBSD 11.0 x64",
		"type": "public"
	}`)

	var s schema.ISO
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	iso := ISOFromSchema(s)

	if iso.ID != 4711 {
		t.Errorf("unexpected ID: %v", iso.ID)
	}
	if iso.Name != "FreeBSD-11.0-RELEASE-amd64-dvd1" {
		t.Errorf("unexpected name: %v", iso.Name)
	}
	if iso.Description != "FreeBSD 11.0 x64" {
		t.Errorf("unexpected description: %v", iso.Description)
	}
	if iso.Type != ISOTypePublic {
		t.Errorf("unexpected type: %v", iso.Type)
	}
}

func TestServerFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "server.example.com",
		"status": "running",
		"created": "2017-08-16T17:29:14+00:00",
		"public_net": {
			"ipv4": {
				"ip": "1.2.3.4",
				"blocked": false,
				"dns_ptr": "server01.example.com"
			},
			"ipv6": {
				"ip": "2a01:4f8:1c11:3400::/64",
				"blocked": false,
				"dns_ptr": [
					{
						"ip": "2a01:4f8:1c11:3400::1/64",
						"dns_ptr": "server01.example.com"
					}
				]
			}
		},
		"server_type": {
			"id": 2
		},
		"outgoing_traffic": 123456,
		"ingoing_traffic": 7891011,
		"included_traffic": 654321,
		"backup_window": "22-02",
		"rescue_enabled": true,
		"iso": {
			"id": 4711,
			"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
			"description": "FreeBSD 11.0 x64",
			"type": "public"
		}
	}`)

	var s schema.Server
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	server := ServerFromSchema(s)

	if server.ID != 1 {
		t.Errorf("unexpected ID: %v", server.ID)
	}
	if server.Name != "server.example.com" {
		t.Errorf("unexpected name: %v", server.Name)
	}
	if server.Status != ServerStatusRunning {
		t.Errorf("unexpected status: %v", server.Status)
	}
	if !server.Created.Equal(time.Date(2017, 8, 16, 17, 29, 14, 0, time.UTC)) {
		t.Errorf("unexpected created date: %v", server.Created)
	}
	if server.PublicNet.IPv4.IP != "1.2.3.4" {
		t.Errorf("unexpected public net IPv4 IP: %v", server.PublicNet.IPv4.IP)
	}
	if server.ServerType.ID != 2 {
		t.Errorf("unexpected server type ID: %v", server.ServerType.ID)
	}
	if server.IncludedTraffic != 654321 {
		t.Errorf("unexpected included traffic: %v", server.IncludedTraffic)
	}
	if server.OutgoingTraffic != 123456 {
		t.Errorf("unexpected outgoing traffic: %v", server.OutgoingTraffic)
	}
	if server.IngoingTraffic != 7891011 {
		t.Errorf("unexpected ingoing traffic: %v", server.IngoingTraffic)
	}
	if server.BackupWindow != "22-02" {
		t.Errorf("unexpected backup window: %v", server.BackupWindow)
	}
	if !server.RescueEnabled {
		t.Errorf("unexpected rescue enabled state: %v", server.RescueEnabled)
	}
	if server.ISO == nil || server.ISO.ID != 4711 {
		t.Errorf("unexpected ISO: %v", server.ISO)
	}
}

func TestServerFromSchemaNoTraffic(t *testing.T) {
	data := []byte(`{
		"outgoing_traffic": null,
		"ingoing_traffic": null
	}`)

	var s schema.Server
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	server := ServerFromSchema(s)

	if server.OutgoingTraffic != 0 {
		t.Errorf("unexpected outgoing traffic: %v", server.OutgoingTraffic)
	}
	if server.IngoingTraffic != 0 {
		t.Errorf("unexpected ingoing traffic: %v", server.IngoingTraffic)
	}
}

func TestServerPublicNetFromSchema(t *testing.T) {
	data := []byte(`{
		"ipv4": {
			"ip": "1.2.3.4",
			"blocked": false,
			"dns_ptr": "server.example.com"
		},
		"ipv6": {
        		"ip": "2a01:4f8:1c19:1403::/64",
        		"blocked": false,
        		"dns_ptr": []
      		},
      		"floating_ips": [4]
	}`)

	var s schema.ServerPublicNet
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	publicNet := ServerPublicNetFromSchema(s)

	if publicNet.IPv4.IP != "1.2.3.4" {
		t.Errorf("unexpected IPv4 IP: %v", publicNet.IPv4.IP)
	}
	if publicNet.IPv6.IP != "2a01:4f8:1c19:1403::/64" {
		t.Errorf("unexpected IPv6 IP: %v", publicNet.IPv6.IP)
	}
	if len(publicNet.FloatingIPs) != 1 || publicNet.FloatingIPs[0].ID != 4 {
		t.Errorf("unexpected Floating IPs: %v", publicNet.FloatingIPs)
	}
}

func TestServerPublicNetIPv4FromSchema(t *testing.T) {
	data := []byte(`{
		"ip": "1.2.3.4",
		"blocked": true,
		"dns_ptr": "server.example.com"
	}`)

	var s schema.ServerPublicNetIPv4
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	ipv4 := ServerPublicNetIPv4FromSchema(s)

	if ipv4.IP != "1.2.3.4" {
		t.Errorf("unexpected IP: %v", ipv4.IP)
	}
	if !ipv4.Blocked {
		t.Errorf("unexpected blocked state: %v", ipv4.Blocked)
	}
	if ipv4.DNSPtr != "server.example.com" {
		t.Errorf("unexpected DNS ptr: %v", ipv4.DNSPtr)
	}
}

func TestServerPublicNetIPv6FromSchema(t *testing.T) {
	data := []byte(`{
		"ip": "2a01:4f8:1c11:3400::/64",
		"blocked": true,
		"dns_ptr": [
			{
				"ip": "2a01:4f8:1c11:3400::1/64",
				"blocked": "server01.example.com"
			}
		]
	}`)

	var s schema.ServerPublicNetIPv6
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	ipv6 := ServerPublicNetIPv6FromSchema(s)

	if ipv6.IP != "2a01:4f8:1c11:3400::/64" {
		t.Errorf("unexpected IP: %v", ipv6.IP)
	}
	if !ipv6.Blocked {
		t.Errorf("unexpected blocked state: %v", ipv6.Blocked)
	}
	if len(ipv6.DNSPtr) != 1 {
		t.Errorf("unexpected DNS ptr: %v", ipv6.DNSPtr)
	}
}

func TestServerPublicNetIPv6DNSPtrFromSchema(t *testing.T) {
	data := []byte(`{
		"ip": "2a01:4f8:1c11:3400::1/64",
		"dns_ptr": "server01.example.com"
	}`)

	var s schema.ServerPublicNetIPv6DNSPtr
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	dnsPtr := ServerPublicNetIPv6DNSPtrFromSchema(s)

	if dnsPtr.IP != "2a01:4f8:1c11:3400::1/64" {
		t.Errorf("unexpected IP: %v", dnsPtr.IP)
	}
	if dnsPtr.DNSPtr != "server01.example.com" {
		t.Errorf("unexpected DNS ptr: %v", dnsPtr.DNSPtr)
	}
}

func TestServerTypeFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "cx10",
		"description": "description",
		"cores": 4,
		"memory": 1.0,
		"disk": 20,
		"storage_type": "local"
	}`)

	var s schema.ServerType
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	serverType := ServerTypeFromSchema(s)

	if serverType.ID != 1 {
		t.Errorf("unexpected ID: %v", serverType.ID)
	}
	if serverType.Name != "cx10" {
		t.Errorf("unexpected name: %q", serverType.Name)
	}
	if serverType.Description != "description" {
		t.Errorf("unexpected description: %q", serverType.Description)
	}
	if serverType.Cores != 4 {
		t.Errorf("unexpected cores: %v", serverType.Cores)
	}
	if serverType.Memory != 1.0 {
		t.Errorf("unexpected memory: %v", serverType.Memory)
	}
	if serverType.Disk != 20 {
		t.Errorf("unexpected disk: %v", serverType.Disk)
	}
	if serverType.StorageType != StorageTypeLocal {
		t.Errorf("unexpected storage type: %q", serverType.StorageType)
	}
}

func TestSSHKeyFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 2323,
		"name": "My key",
		"fingerprint": "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c",
		"public_key": "ssh-rsa AAAjjk76kgf...Xt"
	}`)

	var s schema.SSHKey
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	sshKey := SSHKeyFromSchema(s)

	if sshKey.ID != 2323 {
		t.Errorf("unexpected ID: %v", sshKey.ID)
	}
	if sshKey.Name != "My key" {
		t.Errorf("unexpected name: %v", sshKey.Name)
	}
	if sshKey.Fingerprint != "b7:2f:30:a0:2f:6c:58:6c:21:04:58:61:ba:06:3b:2c" {
		t.Errorf("unexpected fingerprint: %v", sshKey.Fingerprint)
	}
	if sshKey.PublicKey != "ssh-rsa AAAjjk76kgf...Xt" {
		t.Errorf("unexpected public key: %v", sshKey.PublicKey)
	}
}
