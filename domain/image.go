package domain

// Image struct for DB
type Image struct {
	Base
	Path     string  `json:"path,omitempty" binding:"required"`
	Category *string `json:"category" swaggerignore:"true"`
} // @name ImageResponse

// ImageArray struct
type ImageArray struct {
	ID   string `json:"id" gorm:"column:id"`
	Path string `json:"path" gorm:"column:path"`
} //@name Image
