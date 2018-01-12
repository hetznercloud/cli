package schema

// FloatingIP defines the schema of a Floating IP.
type FloatingIP struct {
	ID           int                `json:"id"`
	Description  *string            `json:"description"`
	IP           string             `json:"ip"`
	Type         string             `json:"type"`
	Server       *int               `json:"server"`
	DNSPtr       []FloatingIPDNSPtr `json:"dns_ptr"`
	HomeLocation Location           `json:"home_location"`
	Blocked      bool               `json:"blocked"`
}

// FloatingIPDNSPtr contains reverse DNS information for a
// IPv4 or IPv6 Floating IP.
type FloatingIPDNSPtr struct {
	IP     string `json:"ip"`
	DNSPtr string `json:"dns_ptr"`
}

// FloatingIPGetResponse defines the schema of the response when
// retrieving a single Floating IP.
type FloatingIPGetResponse struct {
	FloatingIP FloatingIP `json:"floating_ip"`
}

// FloatingIPUpdateRequest defines the schema of the request to update a Floating IP.
type FloatingIPUpdateRequest struct {
	Description string `json:"description,omitempty"`
}

// FloatingIPUpdateResponse defines the schema of the response when updating a Floating IP.
type FloatingIPUpdateResponse struct {
	FloatingIP FloatingIP `json:"floating_ip"`
}

// FloatingIPListResponse defines the schema of the response when
// listing Floating IPs.
type FloatingIPListResponse struct {
	FloatingIPs []FloatingIP `json:"floating_ips"`
}

// FloatingIPCreateRequest defines the schema of the request to
// create a Floating IP.
type FloatingIPCreateRequest struct {
	Type         string  `json:"type"`
	HomeLocation *string `json:"home_location,omitempty"`
	Server       *int    `json:"server,omitempty"`
	Description  *string `json:"description,omitempty"`
}

// FloatingIPCreateResponse defines the schema of the response
// when creating a Floating IP.
type FloatingIPCreateResponse struct {
	FloatingIP FloatingIP `json:"floating_ip"`
	Action     *Action    `json:"action"`
}

// FloatingIPActionAssignRequest defines the schema of the request to
// create an assign Floating IP action.
type FloatingIPActionAssignRequest struct {
	Server int `json:"server"`
}

// FloatingIPActionAssignResponse defines the schema of the response when
// creating an assign action.
type FloatingIPActionAssignResponse struct {
	Action Action `json:"action"`
}

// FloatingIPActionUnassignRequest defines the schema of the request to
// create an unassign Floating IP action.
type FloatingIPActionUnassignRequest struct{}

// FloatingIPActionUnassignResponse defines the schema of the response when
// creating an unassign action.
type FloatingIPActionUnassignResponse struct {
	Action Action `json:"action"`
}
