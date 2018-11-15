package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string `gorm:"unique" form:"name"`
	Token string
	// Db    services.DB `gorm:"-"`
}
