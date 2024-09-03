package entity

import (
   "time"
   "gorm.io/gorm"
)

type ShowTimes struct {
   gorm.Model

   Showdate  time.Time 

   //FK
   MovieID *uint
   Movie   Movie `gorm:"foreignKey:MovieID"`

   TheaterID *uint
   Theater  Theater `gorm:"foreignKey:TheaterID"`

   //onetomany
   Tickets []Ticket `gorm:"foreignKey:ShowTimeID"`
}
