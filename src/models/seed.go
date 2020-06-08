package models

// Seed model
type Seed struct {
	Base
	Name string `gorm:"type:varchar(30);unique;not null"`
}