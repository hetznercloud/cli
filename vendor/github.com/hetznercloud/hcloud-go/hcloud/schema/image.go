package schema

import "time"

// Image defines the schema of an image.
type Image struct {
	ID          int               `json:"id"`
	Status      string            `json:"status"`
	Type        string            `json:"type"`
	Name        *string           `json:"name"`
	Description string            `json:"description"`
	ImageSize   *float32          `json:"image_size"`
	DiskSize    float32           `json:"disk_size"`
	Created     time.Time         `json:"created"`
	CreatedFrom *ImageCreatedFrom `json:"created_from"`
	BoundTo     *int              `json:"bound_to"`
	OSFlavor    string            `json:"os_flavor"`
	OSVersion   *string           `json:"os_version"`
	RapidDeploy bool              `json:"rapid_deploy"`
}

// ImageCreatedFrom defines the schema of the images created from reference.
type ImageCreatedFrom struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ImageGetResponse defines the schema of the response when
// retrieving a single image.
type ImageGetResponse struct {
	Image Image `json:"image"`
}

// ImageListResponse defines the schema of the response when
// listing images.
type ImageListResponse struct {
	Images []Image `json:"images"`
}

// ImageUpdateRequest defines the schema of the request to update an image.
type ImageUpdateRequest struct {
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"`
}

// ImageUpdateResponse defines the schema of the response when updating an image.
type ImageUpdateResponse struct {
	Image Image `json:"image"`
}
