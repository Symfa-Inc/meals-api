package domain

// User model
type User struct {
	Base
	FirstName   string  `gorm:"type:varchar(20)" json:"firstName"`
	LastName    string  `gorm:"type:varchar(20)" json:"lastName"`
	Email       string  `gorm:"type:varchar(30);not null" json:"email"`
	Password    string  `gorm:"type:varchar(100);not null" json:"-"`
	Role        string  `sql:"type:user_roles" json:"role"`
	CompanyType *string `sql:"type:company_types" gorm:"type:varchar(20);null" json:"companyType"`
	Status      *string `sql:"type:status_types" gorm:"type:varchar(10);null" json:"status"`
}
