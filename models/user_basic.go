package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity     string `gorm:"column:identity;type:varchar(36);" json:"identity"`
	Name         string `gorm:"column:name;type:varchar(100);" json:"name"`
	Password     string `gorm:"column:password;type:varchar(32);" json:"password"`
	Phone        string `gorm:"column:phone;type:varchar(20);" json:"phone"`
	Email        string `gorm:"column:email;type:varchar(100);" json:"email"`
	AcNumber     int64  `gorm:"column:ac_number;type:int(11);" json:"ac_number"`
	SubmitNumber int64  `gorm:"column:submit_number;type:int(11);" json:"submit_number"`
	IsAdmin      int    `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"` //[0: false, 1: true]
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
