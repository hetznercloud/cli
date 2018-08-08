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
	if len(action.Resources) == 1 {
		if action.Resources[0].ID != 42 {
			t.Errorf("unexpected id in resources[0].ID: %v", action.Resources[0].ID)
		}
		if action.Resources[0].Type != ActionResourceTypeServer {
			t.Errorf("unexpected type in resources[0].Type: %v", action.Resources[0].Type)
		}
	} else {
		t.Errorf("unexpected number of resources")
	}
}

func TestFloatingIPFromSchema(t *testing.T) {
	t.Run("IPv6", func(t *testing.T) {
		data := []byte(`{
			"id": 4711,
			"description": "Web Frontend",
			"ip": "2001:db8::/64",
			"type": "ipv6",
			"server": null,
			"dns_ptr": [],
			"blocked": true,
			"home_location": {
				"id": 1,
				"name": "fsn1",
				"description": "Falkenstein DC Park 1",
				"country": "DE",
				"city": "Falkenstein",
				"latitude": 50.47612,
				"longitude": 12.370071
			},
			"protection": {
				"delete": true
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
		if !floatingIP.Blocked {
			t.Errorf("unexpected value for Blocked: %v", floatingIP.Blocked)
		}
		if floatingIP.Description != "Web Frontend" {
			t.Errorf("unexpected description: %v", floatingIP.Description)
		}
		if floatingIP.IP.String() != "2001:db8::" {
			t.Errorf("unexpected IP: %v", floatingIP.IP)
		}
		if floatingIP.Type != FloatingIPTypeIPv6 {
			t.Errorf("unexpected Type: %v", floatingIP.Type)
		}
		if floatingIP.Server != nil {
			t.Errorf("unexpected Server: %v", floatingIP.Server)
		}
		if floatingIP.DNSPtr == nil || floatingIP.DNSPtrForIP(floatingIP.IP) != "" {
			t.Errorf("unexpected DNS ptr: %v", floatingIP.DNSPtr)
		}
		if floatingIP.HomeLocation == nil || floatingIP.HomeLocation.ID != 1 {
			t.Errorf("unexpected home location: %v", floatingIP.HomeLocation)
		}
		if !floatingIP.Protection.Delete {
			t.Errorf("unexpected Protection.Delete: %v", floatingIP.Protection.Delete)
		}
	})

	t.Run("IPv4", func(t *testing.T) {
		data := []byte(`{
			"id": 4711,
			"description": "Web Frontend",
			"ip": "131.232.99.1",
			"type": "ipv4",
			"server": 42,
			"dns_ptr": [{
				"ip": "131.232.99.1",
				"dns_ptr": "fip01.example.com"
			}],
			"blocked": false,
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
		if floatingIP.Blocked {
			t.Errorf("unexpected value for Blocked: %v", floatingIP.Blocked)
		}
		if floatingIP.Description != "Web Frontend" {
			t.Errorf("unexpected description: %v", floatingIP.Description)
		}
		if floatingIP.IP.String() != "131.232.99.1" {
			t.Errorf("unexpected IP: %v", floatingIP.IP)
		}
		if floatingIP.Type != FloatingIPTypeIPv4 {
			t.Errorf("unexpected type: %v", floatingIP.Type)
		}
		if floatingIP.Server == nil || floatingIP.Server.ID != 42 {
			t.Errorf("unexpected server: %v", floatingIP.Server)
		}
		if floatingIP.DNSPtr == nil || floatingIP.DNSPtrForIP(floatingIP.IP) != "fip01.example.com" {
			t.Errorf("unexpected DNS ptr: %v", floatingIP.DNSPtr)
		}
		if floatingIP.HomeLocation == nil || floatingIP.HomeLocation.ID != 1 {
			t.Errorf("unexpected home location: %v", floatingIP.HomeLocation)
		}
	})
}

func TestISOFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
		"description": "FreeBSD 11.0 x64",
		"type": "public",
		"deprecated": "2018-02-28T00:00:00+00:00"
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
	if iso.Deprecated.IsZero() {
		t.Errorf("unexpected value for deprecated: %v", iso.Deprecated)
	}
}

func TestDatacenterFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 1,
		"name": "fsn1-dc8",
		"description": "Falkenstein 1 DC 8",
		"location": {
			"id": 1,
			"name": "fsn1",
			"description": "Falkenstein DC Park 1",
			"country": "DE",
			"city": "Falkenstein",
			"latitude": 50.47612,
			"longitude": 12.370071
		},
		"server_types": {
			"supported": [
				1,
				1,
				2,
				3
			],
			"available": [
				1,
				1,
				2,
				3
			]
		}
	}`)

	var s schema.Datacenter
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	datacenter := DatacenterFromSchema(s)
	if datacenter.ID != 1 {
		t.Errorf("unexpected ID: %v", datacenter.ID)
	}
	if datacenter.Name != "fsn1-dc8" {
		t.Errorf("unexpected Name: %v", datacenter.Name)
	}
	if datacenter.Location == nil || datacenter.Location.ID != 1 {
		t.Errorf("unexpected Location: %v", datacenter.Location)
	}
	if len(datacenter.ServerTypes.Available) != 4 {
		t.Errorf("unexpected ServerTypes.Available (should be 4): %v", len(datacenter.ServerTypes.Available))
	}
	if len(datacenter.ServerTypes.Supported) != 4 {
		t.Errorf("unexpected ServerTypes.Supported length (should be 4): %v", len(datacenter.ServerTypes.Supported))
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
		"image": {
			"id": 4711,
			"type": "system",
			"status": "available",
			"name": "ubuntu16.04-standard-x64",
			"description": "Ubuntu 16.04 Standard 64 bit",
			"image_size": 2.3,
			"disk_size": 10,
			"created": "2017-08-16T17:29:14+00:00",
			"created_from": {
				"id": 1,
				"name": "Server"
			},
			"bound_to": 1,
			"os_flavor": "ubuntu",
			"os_version": "16.04",
			"rapid_deploy": false
		},
		"iso": {
			"id": 4711,
			"name": "FreeBSD-11.0-RELEASE-amd64-dvd1",
			"description": "FreeBSD 11.0 x64",
			"type": "public"
		},
		"datacenter": {
			"id": 1,
			"name": "fsn1-dc8",
			"description": "Falkenstein 1 DC 8",
			"location": {
				"id": 1,
				"name": "fsn1",
				"description": "Falkenstein DC Park 1",
				"country": "DE",
				"city": "Falkenstein",
				"latitude": 50.47612,
				"longitude": 12.370071
			}
		},
		"protection": {
			"delete": true,
			"rebuild": true
		},
		"locked": true
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
	if server.PublicNet.IPv4.IP.String() != "1.2.3.4" {
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
	if server.Image == nil || server.Image.ID != 4711 {
		t.Errorf("unexpected Image: %v", server.Image)
	}
	if server.ISO == nil || server.ISO.ID != 4711 {
		t.Errorf("unexpected ISO: %v", server.ISO)
	}
	if server.Datacenter == nil || server.Datacenter.ID != 1 {
		t.Errorf("unexpected Datacenter: %v", server.Datacenter)
	}
	if !server.Locked {
		t.Errorf("unexpected value for Locked: %v", server.Locked)
	}
	if !server.Protection.Delete {
		t.Errorf("unexpected value for Protection.Delete: %v", server.Protection.Delete)
	}
	if !server.Protection.Rebuild {
		t.Errorf("unexpected value for Protection.Rebuild: %v", server.Protection.Rebuild)
	}
}

func TestServerFromSchemaNoTraffic(t *testing.T) {
	data := []byte(`{
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

	if publicNet.IPv4.IP.String() != "1.2.3.4" {
		t.Errorf("unexpected IPv4 IP: %v", publicNet.IPv4.IP)
	}
	if publicNet.IPv6.Network.String() != "2a01:4f8:1c19:1403::/64" {
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

	if ipv4.IP.String() != "1.2.3.4" {
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

	if ipv6.Network.String() != "2a01:4f8:1c11:3400::/64" {
		t.Errorf("unexpected IP: %v", ipv6.IP)
	}
	if !ipv6.Blocked {
		t.Errorf("unexpected blocked state: %v", ipv6.Blocked)
	}
	if len(ipv6.DNSPtr) != 1 {
		t.Errorf("unexpected DNS ptr: %v", ipv6.DNSPtr)
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
		"storage_type": "local",
		"cpu_type": "shared",
		"prices": [
			{
				"location": "fsn1",
				"price_hourly": {
					"net": "1",
					"gross": "1.19"
				},
				"price_monthly": {
					"net": "1",
					"gross": "1.19"
				}
			}
		]
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
	if serverType.CPUType != CPUTypeShared {
		t.Errorf("unexpected cpu type: %q", serverType.CPUType)
	}
	if len(serverType.Pricings) != 1 {
		t.Errorf("unexpected number of pricings: %d", len(serverType.Pricings))
	} else {
		if serverType.Pricings[0].Location.Name != "fsn1" {
			t.Errorf("unexpected location name: %v", serverType.Pricings[0].Location.Name)
		}
		if serverType.Pricings[0].Hourly.Net != "1" {
			t.Errorf("unexpected hourly net price: %v", serverType.Pricings[0].Hourly.Net)
		}
		if serverType.Pricings[0].Hourly.Gross != "1.19" {
			t.Errorf("unexpected hourly gross price: %v", serverType.Pricings[0].Hourly.Gross)
		}
		if serverType.Pricings[0].Monthly.Net != "1" {
			t.Errorf("unexpected monthly net price: %v", serverType.Pricings[0].Monthly.Net)
		}
		if serverType.Pricings[0].Monthly.Gross != "1.19" {
			t.Errorf("unexpected monthly gross price: %v", serverType.Pricings[0].Monthly.Gross)
		}
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

func TestErrorFromSchema(t *testing.T) {
	t.Run("service_error", func(t *testing.T) {
		data := []byte(`{
			"code": "service_error",
			"message": "An error occured",
			"details": {}
		}`)

		var s schema.Error
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		err := ErrorFromSchema(s)

		if err.Code != "service_error" {
			t.Errorf("unexpected code: %v", err.Code)
		}
		if err.Message != "An error occured" {
			t.Errorf("unexpected message: %v", err.Message)
		}
	})

	t.Run("invalid_input", func(t *testing.T) {
		data := []byte(`{
			"code": "invalid_input",
			"message": "invalid input",
			"details": {
				"fields": [
					{
						"name": "broken_field",
						"messages": ["is required"]
					}
				]
			}
		}`)

		var s schema.Error
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatal(err)
		}
		err := ErrorFromSchema(s)

		if err.Code != "invalid_input" {
			t.Errorf("unexpected Code: %v", err.Code)
		}
		if err.Message != "invalid input" {
			t.Errorf("unexpected Message: %v", err.Message)
		}
		if d, ok := err.Details.(ErrorDetailsInvalidInput); !ok {
			t.Fatalf("unexpected Details type (should be ErrorDetailsInvalidInput): %v", err.Details)
		} else {
			if len(d.Fields) != 1 {
				t.Fatalf("unexpected Details.Fields length (should be 1): %v", d.Fields)
			}
			if d.Fields[0].Name != "broken_field" {
				t.Errorf("unexpected Details.Fields[0].Name: %v", d.Fields[0].Name)
			}
			if len(d.Fields[0].Messages) != 1 {
				t.Fatalf("unexpected Details.Fields[0].Messages length (should be 1): %v", d.Fields[0].Messages)
			}
			if d.Fields[0].Messages[0] != "is required" {
				t.Errorf("unexpected Details.Fields[0].Messages[0]: %v", d.Fields[0].Messages[0])
			}
		}
	})
}

func TestPaginationFromSchema(t *testing.T) {
	data := []byte(`{
		"page": 2,
		"per_page": 25,
		"previous_page": 1,
		"next_page": 3,
		"last_page": 13,
		"total_entries": 322
	}`)

	var s schema.MetaPagination
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	p := PaginationFromSchema(s)

	if p.Page != 2 {
		t.Errorf("unexpected page: %v", p.Page)
	}
	if p.PerPage != 25 {
		t.Errorf("unexpected per page: %v", p.PerPage)
	}
	if p.PreviousPage != 1 {
		t.Errorf("unexpected previous page: %v", p.PreviousPage)
	}
	if p.NextPage != 3 {
		t.Errorf("unexpected next page: %d", p.NextPage)
	}
	if p.LastPage != 13 {
		t.Errorf("unexpected last page: %d", p.LastPage)
	}
	if p.TotalEntries != 322 {
		t.Errorf("unexpected total entries: %d", p.TotalEntries)
	}
}

func TestImageFromSchema(t *testing.T) {
	data := []byte(`{
		"id": 4711,
		"type": "system",
		"status": "available",
		"name": "ubuntu16.04-standard-x64",
		"description": "Ubuntu 16.04 Standard 64 bit",
		"image_size": 2.3,
		"disk_size": 10,
		"created": "2016-01-30T23:55:01Z",
		"created_from": {
			"id": 1,
			"name": "my-server1"
		},
		"bound_to": 1,
		"os_flavor": "ubuntu",
		"os_version": "16.04",
		"rapid_deploy": false,
		"protection": {
			"delete": true
		},
		"deprecated": "2018-02-28T00:00:00+00:00"
	}`)

	var s schema.Image
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	image := ImageFromSchema(s)

	if image.ID != 4711 {
		t.Errorf("unexpected ID: %v", image.ID)
	}
	if image.Type != ImageTypeSystem {
		t.Errorf("unexpected Type: %v", image.Type)
	}
	if image.Status != ImageStatusAvailable {
		t.Errorf("unexpected Status: %v", image.Status)
	}
	if image.Name != "ubuntu16.04-standard-x64" {
		t.Errorf("unexpected Name: %v", image.Name)
	}
	if image.Description != "Ubuntu 16.04 Standard 64 bit" {
		t.Errorf("unexpected Description: %v", image.Description)
	}
	if image.ImageSize != 2.3 {
		t.Errorf("unexpected ImageSize: %v", image.ImageSize)
	}
	if image.DiskSize != 10 {
		t.Errorf("unexpected DiskSize: %v", image.DiskSize)
	}
	if !image.Created.Equal(time.Date(2016, 1, 30, 23, 55, 1, 0, time.UTC)) {
		t.Errorf("unexpected Created: %v", image.Created)
	}
	if image.CreatedFrom == nil || image.CreatedFrom.ID != 1 || image.CreatedFrom.Name != "my-server1" {
		t.Errorf("unexpected CreatedFrom: %v", image.CreatedFrom)
	}
	if image.BoundTo == nil || image.BoundTo.ID != 1 {
		t.Errorf("unexpected BoundTo: %v", image.BoundTo)
	}
	if image.OSVersion != "16.04" {
		t.Errorf("unexpected OSVersion: %v", image.OSVersion)
	}
	if image.OSFlavor != "ubuntu" {
		t.Errorf("unexpected OSFlavor: %v", image.OSFlavor)
	}
	if image.RapidDeploy {
		t.Errorf("unexpected RapidDeploy: %v", image.RapidDeploy)
	}
	if !image.Protection.Delete {
		t.Errorf("unexpected Protection.Delete: %v", image.Protection.Delete)
	}
	if image.Deprecated.IsZero() {
		t.Errorf("unexpected value for Deprecated: %v", image.Deprecated)
	}
}

func TestPricingFromSchema(t *testing.T) {
	data := []byte(`{
		"currency": "EUR",
		"vat_rate": "19.00",
		"image": {
			"price_per_gb_month": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"floating_ip": {
			"price_monthly": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"traffic": {
			"price_per_tb": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"server_backup": {
			"percentage": "20"
		},
		"server_types": [
			{
				"id": 4,
				"name": "CX11",
				"prices": [
					{
						"location": "fsn1",
						"price_hourly": {
							"net": "1",
							"gross": "1.19"
						},
						"price_monthly": {
							"net": "1",
							"gross": "1.19"
						}
					}
				]
			}
		]
	}`)

	var s schema.Pricing
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	pricing := PricingFromSchema(s)

	if pricing.Image.PerGBMonth.Currency != "EUR" {
		t.Errorf("unexpected Image.PerGBMonth.Currency: %v", pricing.Image.PerGBMonth.Currency)
	}
	if pricing.Image.PerGBMonth.VATRate != "19.00" {
		t.Errorf("unexpected Image.PerGBMonth.VATRate: %v", pricing.Image.PerGBMonth.VATRate)
	}
	if pricing.Image.PerGBMonth.Net != "1" {
		t.Errorf("unexpected Image.PerGBMonth.Net: %v", pricing.Image.PerGBMonth.Net)
	}
	if pricing.Image.PerGBMonth.Gross != "1.19" {
		t.Errorf("unexpected Image.PerGBMonth.Gross: %v", pricing.Image.PerGBMonth.Gross)
	}

	if pricing.FloatingIP.Monthly.Currency != "EUR" {
		t.Errorf("unexpected FloatingIP.Monthly.Currency: %v", pricing.FloatingIP.Monthly.Currency)
	}
	if pricing.FloatingIP.Monthly.VATRate != "19.00" {
		t.Errorf("unexpected FloatingIP.Monthly.VATRate: %v", pricing.FloatingIP.Monthly.VATRate)
	}
	if pricing.FloatingIP.Monthly.Net != "1" {
		t.Errorf("unexpected FloatingIP.Monthly.Net: %v", pricing.FloatingIP.Monthly.Net)
	}
	if pricing.FloatingIP.Monthly.Gross != "1.19" {
		t.Errorf("unexpected FloatingIP.Monthly.Gross: %v", pricing.FloatingIP.Monthly.Gross)
	}

	if pricing.Traffic.PerTB.Currency != "EUR" {
		t.Errorf("unexpected Traffic.PerTB.Currency: %v", pricing.Traffic.PerTB.Currency)
	}
	if pricing.Traffic.PerTB.VATRate != "19.00" {
		t.Errorf("unexpected Traffic.PerTB.VATRate: %v", pricing.Traffic.PerTB.VATRate)
	}
	if pricing.Traffic.PerTB.Net != "1" {
		t.Errorf("unexpected Traffic.PerTB.Net: %v", pricing.Traffic.PerTB.Net)
	}
	if pricing.Traffic.PerTB.Gross != "1.19" {
		t.Errorf("unexpected Traffic.PerTB.Gross: %v", pricing.Traffic.PerTB.Gross)
	}

	if pricing.ServerBackup.Percentage != "20" {
		t.Errorf("unexpected ServerBackup.Percentage: %v", pricing.ServerBackup.Percentage)
	}

	if len(pricing.ServerTypes) != 1 {
		t.Errorf("unexpected number of server types: %d", len(pricing.ServerTypes))
	} else {
		p := pricing.ServerTypes[0]

		if p.ServerType.ID != 4 {
			t.Errorf("unexpected ServerType.ID: %d", p.ServerType.ID)
		}
		if p.ServerType.Name != "CX11" {
			t.Errorf("unexpected ServerType.Name: %v", p.ServerType.Name)
		}

		if len(p.Pricings) != 1 {
			t.Errorf("unexpected number of prices: %d", len(p.Pricings))
		} else {
			if p.Pricings[0].Location.Name != "fsn1" {
				t.Errorf("unexpected Location.Name: %v", p.Pricings[0].Location.Name)
			}

			if p.Pricings[0].Hourly.Currency != "EUR" {
				t.Errorf("unexpected Hourly.Currency: %v", p.Pricings[0].Hourly.Currency)
			}
			if p.Pricings[0].Hourly.VATRate != "19.00" {
				t.Errorf("unexpected Hourly.VATRate: %v", p.Pricings[0].Hourly.VATRate)
			}
			if p.Pricings[0].Hourly.Net != "1" {
				t.Errorf("unexpected Hourly.Net: %v", p.Pricings[0].Hourly.Net)
			}
			if p.Pricings[0].Hourly.Gross != "1.19" {
				t.Errorf("unexpected Hourly.Gross: %v", p.Pricings[0].Hourly.Gross)
			}

			if p.Pricings[0].Monthly.Currency != "EUR" {
				t.Errorf("unexpected Monthly.Currency: %v", p.Pricings[0].Monthly.Currency)
			}
			if p.Pricings[0].Monthly.VATRate != "19.00" {
				t.Errorf("unexpected Monthly.VATRate: %v", p.Pricings[0].Monthly.VATRate)
			}
			if p.Pricings[0].Monthly.Net != "1" {
				t.Errorf("unexpected Monthly.Net: %v", p.Pricings[0].Monthly.Net)
			}
			if p.Pricings[0].Monthly.Gross != "1.19" {
				t.Errorf("unexpected Monthly.Gross: %v", p.Pricings[0].Monthly.Gross)
			}
		}
	}
}
