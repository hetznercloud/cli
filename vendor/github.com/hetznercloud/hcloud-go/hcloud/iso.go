package hcloud

// ISO represents an ISO image in the Hetzner Cloud.
type ISO struct {
	ID          int
	Name        string
	Description string
	Type        ISOType
}

// ISOType specifies the type of an ISO image.
type ISOType string

const (
	// ISOTypePublic is the type of a public ISO image.
	ISOTypePublic ISOType = "public"

	// ISOTypePrivate is the type of a private ISO image.
	ISOTypePrivate = "private"
)
