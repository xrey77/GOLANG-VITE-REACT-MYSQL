package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          int       `gorm:"type:integer;primarykey"`
	Lastname    string    `gorm:"size:255;default:null"`
	Firstname   string    `gorm:"size:255;default:null"`
	Email       string    `gorm:"size:255;uniqueIndex;not:null"`
	Mobile      string    `gorm:"size:255;default:null"`
	Username    string    `gorm:"size:255;uniqueIndex;not:null"`
	Password    string    `gorm:"size:255"`
	Isactivated int       `gorm:"type:integer;default:1"`
	Isblocked   int       `gorm:"type:integer;default:0"`
	Userpicture string    `gorm:"default:http://127.0.0.1:5000/assets/users/pix.png"`
	Mailtoken   string    `gorm:"type:integer;default:0"`
	Secret      string    `gorm:"type:text;default:null"`
	Qrcodeurl   string    `gorm:"type:text;default:null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)"`
	Role_id     int       `gorm:"type:integer"`

	Roles []Role `gorm:"many2many:user_roles;"`
}
