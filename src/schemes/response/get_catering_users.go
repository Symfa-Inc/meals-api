package response

import uuid "github.com/satori/go.uuid"

type GetCateringUser struct {
	ID        uuid.UUID `gorm:"type:uuid;" json:"id"`
	FirstName string    `gorm:"type:varchar(20)" json:"firstName"`
	LastName  string    `gorm:"type:varchar(20)" json:"lastName"`
	Email     string    `gorm:"type:varchar(30);unique;not null" json:"email"`
	Role      string    `sql:"type:user_roles" json:"role"`
	Status    *string   `sql:"type:status_types" gorm:"type:varchar(10);null" json:"status"`
}
