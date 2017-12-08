package schema

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestFloatingIPDNSPtrMarshalJSON(t *testing.T) {
	t.Run("IPv4", func(t *testing.T) {
		s := "foo.example"
		data, err := json.Marshal(FloatingIPDNSPtr{IPv4: &s})
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, []byte(`"foo.example"`)) {
			t.Errorf("unexpected data: %s", data)
		}
	})
	t.Run("IPv6", func(t *testing.T) {
		data, err := json.Marshal(FloatingIPDNSPtr{
			IPv6: []FloatingIPDNSPtrIPv6{
				{IP: "::1", DNSPtr: "foo.example.com"},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, []byte(`[{"ip":"::1","dns_ptr":"foo.example.com"}]`)) {
			t.Errorf("unexpected data: %s", data)
		}
	})
}

func TestFloatingIPDNSPtrUnmarshalJSON(t *testing.T) {
	t.Run("IPv4", func(t *testing.T) {
		data := []byte(`"foo.example"`)
		var p FloatingIPDNSPtr
		if err := json.Unmarshal(data, &p); err != nil {
			t.Fatal(err)
		}
		if p.IPv4 == nil {
			t.Fatalf("missing IPv4")
		}
		if *p.IPv4 != "foo.example" {
			t.Errorf("unexpected IPv4: %v", *p.IPv4)
		}
		if p.IPv6 != nil {
			t.Error("IPv6 should be nil")
		}
	})
	t.Run("IPv6", func(t *testing.T) {
		data := []byte(`[{"ip":"::1","dns_ptr":"foo.example.com"}]`)
		var p FloatingIPDNSPtr
		if err := json.Unmarshal(data, &p); err != nil {
			t.Fatal(err)
		}
		if p.IPv4 != nil {
			t.Error("IPv4 should be nil")
		}
		if len(p.IPv6) != 1 {
			t.Fatalf("unexpected one entry")
		}
		if p.IPv6[0].IP != "::1" {
			t.Errorf("unexpected IP: %v", p.IPv6[0].IP)
		}
		if p.IPv6[0].DNSPtr != "foo.example.com" {
			t.Errorf("unexpected ptr: %v", p.IPv6[0].DNSPtr)
		}
	})
}
