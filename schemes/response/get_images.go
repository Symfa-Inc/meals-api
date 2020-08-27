package response

// GetImages response scheme
type GetImages struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Category string `json:"category,omitempty"`
} //@name GetImagesResponse
