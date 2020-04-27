package seeds

import "github.com/jinzhu/gorm"

// Seed model
type Seed struct {
	gorm.Model
	Name string `gorm:"type:varchar(30);unique;not null"`
}