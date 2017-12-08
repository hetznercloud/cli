package schema

import "time"

// Server defines the schema of a server.
type Server struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Status          string          `json:"status"`
	Created         time.Time       `json:"created"`
	PublicNet       ServerPublicNet `json:"public_net"`
	ServerType      ServerType      `json:"server_type"`
	IncludedTraffic uint64          `json:"included_traffic"`
	OutgoingTraffic uint64          `json:"outgoing_traffic"`
	IngoingTraffic  uint64          `json:"ingoing_traffic"`
	BackupWindow    string          `json:"backup_window"`
	RescueEnabled   bool            `json:"rescue_enabled"`
	ISO             *ISO            `json:"iso"`
}

// ServerPublicNet defines the schema of a server's
// public network information.
type ServerPublicNet struct {
	IPv4        ServerPublicNetIPv4 `json:"ipv4"`
	IPv6        ServerPublicNetIPv6 `json:"ipv6"`
	FloatingIPs []int               `json:"floating_ips"`
}

// ServerPublicNetIPv4 defines the schema of a server's public
// network information for an IPv4.
type ServerPublicNetIPv4 struct {
	IP      string `json:"ip"`
	Blocked bool   `json:"blocked"`
	DNSPtr  string `json:"dns_ptr"`
}

// ServerPublicNetIPv6 defines the schema of a server's public
// network information for an IPv6.
type ServerPublicNetIPv6 struct {
	IP      string                      `json:"ip"`
	Blocked bool                        `json:"blocked"`
	DNSPtr  []ServerPublicNetIPv6DNSPtr `json:"dns_ptr"`
}

// ServerPublicNetIPv6DNSPtr defines the schema of a server's
// public network information for an IPv6 reverse DNS.
type ServerPublicNetIPv6DNSPtr struct {
	IP     string `json:"ip"`
	DNSPtr string `json:"dns_ptr"`
}

// ServerGetResponse defines the schema of the response when
// retrieving a single server.
type ServerGetResponse struct {
	Server Server `json:"server"`
}

// ServerListResponse defines the schema of the response when
// listing servers.
type ServerListResponse struct {
	Servers []Server `json:"servers"`
}

// ServerCreateRequest defines the schema for the request to
// create a server.
type ServerCreateRequest struct {
	Name       string      `json:"name"`
	ServerType interface{} `json:"server_type"` // int or string
	Image      interface{} `json:"image"`       // int or string
}

// ServerCreateResponse defines the schema of the response when
// creating a server.
type ServerCreateResponse struct {
	Server Server `json:"server"`
	Action Action `json:"action"`
}

// ServerActionPoweronRequest defines the schema for the request to
// create a poweron server action.
type ServerActionPoweronRequest struct{}

// ServerActionPoweronResponse defines the schema of the response when
// creating a poweron server action.
type ServerActionPoweronResponse struct {
	Action Action `json:"action"`
}

// ServerActionPoweroffRequest defines the schema for the request to
// create a poweroff server action.
type ServerActionPoweroffRequest struct{}

// ServerActionPoweroffResponse defines the schema of the response when
// creating a poweroff server action.
type ServerActionPoweroffResponse struct {
	Action Action `json:"action"`
}

// ServerActionRebootRequest defines the schema for the request to
// create a reboot server action.
type ServerActionRebootRequest struct{}

// ServerActionRebootResponse defines the schema of the response when
// creating a reboot server action.
type ServerActionRebootResponse struct {
	Action Action `json:"action"`
}

// ServerActionResetRequest defines the schema for the request to
// create a reset server action.
type ServerActionResetRequest struct{}

// ServerActionResetResponse defines the schema of the response when
// creating a reset server action.
type ServerActionResetResponse struct {
	Action Action `json:"action"`
}

// ServerActionShutdownRequest defines the schema for the request to
// create a shutdown server action.
type ServerActionShutdownRequest struct{}

// ServerActionShutdownResponse defines the schema of the response when
// creating a shutdown server action.
type ServerActionShutdownResponse struct {
	Action Action `json:"action"`
}

// ServerActionResetPasswordRequest defines the schema for the request to
// create a reset_password server action.
type ServerActionResetPasswordRequest struct{}

// ServerActionResetPasswordResponse defines the schema of the response when
// creating a reset_password server action.
type ServerActionResetPasswordResponse struct {
	Action       Action `json:"action"`
	RootPassword string `json:"root_password"`
}
