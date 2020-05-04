package seeds

import "go_api/src/models"

// Seed model
type Seed struct {
	models.Base
	Name string `gorm:"type:varchar(30);unique;not null"`
}