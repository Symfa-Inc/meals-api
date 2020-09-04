package domain

// Seed model
type Seed struct {
	Base
	Name string `gorm:"api_types:varchar(30);unique;not null"`
}
