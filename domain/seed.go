package domain

// Seed model
type Seed struct {
	Base
	Name string `gorm:"varchar(30);unique;not null"`
}
