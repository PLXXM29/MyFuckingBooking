package entity

import (
	

	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	UserName  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	TotalPoint int

	// GenderID ทำหน้าที่เป็น FK
	GenderID uint
	Gender   Gender `gorm:"foriegnKey:GenderID"`
}
