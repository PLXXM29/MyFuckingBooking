package entity

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	BookingTime string
	Status      string

	//FK
	MemberID *uint
	Member   Member `gorm:"foreignKey:MemberID"`

	ShowTimeID *uint
	ShowTime   ShowTimes `gorm:"foreignKey:ShowTimeID"`

	SeatID *uint
	Seat   Seat `gorm:"foreignKey:SeatID"`
}
