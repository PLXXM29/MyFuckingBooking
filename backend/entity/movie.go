package entity

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	MovieName   string
	MovieType   string
	MovieDuration	int
	Director    string
	Actor       string
	Synopsis    string
	ReleaseDate time.Time

	//onetomany   
	ShowTimes []ShowTimes `gorm:"foreignKey:MovieID"`
}
