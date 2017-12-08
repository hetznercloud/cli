package hcloud

// ServerType represents a server type in the Hetzner Cloud.
type ServerType struct {
	ID          int
	Name        string
	Description string
	Cores       int
	Memory      float32
	Disk        int
	StorageType StorageType
}

// StorageType specifies the type of storage.
type StorageType string

const (
	// StorageTypeLocal is the type for local storage.
	StorageTypeLocal StorageType = "local"

	// StorageTypeCeph is the type for remote storage.
	StorageTypeCeph = "ceph"
)
