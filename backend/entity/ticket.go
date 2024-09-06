package entity

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Point int

	//FK
	ShowTimeID *uint
	ShowTime   ShowTimes `gorm:"foreignKey:ShowTimeID"`

	MemberID *uint
	Member   Member `gorm:"foreignKey:MemberID"`

	PaymentID *uint
	Payment   Payment `gorm:"foreignKey:PaymentID"`

	//onetomany
	Seats []Seat `gorm:"foreignKey:SeatID"`
}
