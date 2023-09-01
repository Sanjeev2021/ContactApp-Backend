package contact

import (
	"github.com/jinzhu/gorm"

	"contactapp/models/user"
)

type Contact struct {
	gorm.Model
	User     user.User `gorm:"foreignkey:UserID"`
	UserId   uint
	FullName string `json:"FullName" gotm:"type:varchar(100);not null"`
}
