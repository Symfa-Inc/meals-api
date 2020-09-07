package domain

// Seed model
type Seed struct {
	Base
	Name string `gorm:"url:varchar(30);unique;not null"`
}
