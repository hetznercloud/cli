package hcloud

// Location represents a location in the Hetzner Cloud.
type Location struct {
	ID          int
	Name        string
	Description string
	Country     string
	City        string
	Latitude    float64
	Longitude   float64
}
