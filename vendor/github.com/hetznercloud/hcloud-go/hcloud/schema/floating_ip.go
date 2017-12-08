package schema

import (
	"encoding/json"
	"errors"
)

// FloatingIP defines the schema of a Floating IP.
type FloatingIP struct {
	ID           int              `json:"id"`
	Description  *string          `json:"description"`
	IP           string           `json:"ip"`
	Type         string           `json:"type"`
	Server       *int             `json:"server"`
	DNSPtr       FloatingIPDNSPtr `json:"dns_ptr"`
	HomeLocation Location         `json:"home_location"`
}

// FloatingIPDNSPtr contains reverse DNS information for a
// IPv4 or IPv6 Floating IP.
type FloatingIPDNSPtr struct {
	IPv4 *string
	IPv6 []FloatingIPDNSPtrIPv6
}

// FloatingIPDNSPtrIPv6 defines the schema of reverse DNS
// information for a single IPv6.
type FloatingIPDNSPtrIPv6 struct {
	IP     string `json:"ip"`
	DNSPtr string `json:"dns_ptr"`
}

// MarshalJSON implements json.Marshaler.
func (p FloatingIPDNSPtr) MarshalJSON() ([]byte, error) {
	if p.IPv4 != nil {
		return json.Marshal(p.IPv4)
	}
	if p.IPv6 != nil {
		return json.Marshal(p.IPv6)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *FloatingIPDNSPtr) UnmarshalJSON(data []byte) error {
	var ipv4 string
	if err := json.Unmarshal(data, &ipv4); err == nil {
		p.IPv4 = &ipv4
		return nil
	}

	var ipv6 []FloatingIPDNSPtrIPv6
	if err := json.Unmarshal(data, &ipv6); err == nil {
		p.IPv6 = ipv6
		return nil
	}

	return errors.New("schema: unable to unmarshal dns_ptr")
}

// FloatingIPGetResponse defines the schema of the response when
// retrieving a single Floating IP.
type FloatingIPGetResponse struct {
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
