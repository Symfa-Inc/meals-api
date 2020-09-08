package domain

// Catering model
type Catering struct {
	Base
	Name string `gorm:"type:varchar(30);not null" json:"name,omitempty" binding:"required"`
} //@name CateringsResponse
