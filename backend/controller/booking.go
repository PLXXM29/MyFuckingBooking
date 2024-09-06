package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tanapon395/sa-67-example/config"
	"github.com/tanapon395/sa-67-example/entity"
)

// CreateBooking รับการจองที่นั่ง
func CreateBooking(c *gin.Context) {
	var bookingRequest struct {
		MemberID   uint   `json:"member_id"`
		ShowTimeID uint   `json:"showtime_id"`
		SeatID     []uint `json:"seat_id"`
	}

	// ตรวจสอบว่า JSON ที่ส่งมาถูกต้องหรือไม่
	if err := c.ShouldBindJSON(&bookingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()

	// ตรวจสอบว่าที่นั่งเหล่านี้ว่างอยู่หรือไม่สำหรับรอบเวลานั้น
	for _, seatID := range bookingRequest.SeatID {
		var existingBooking entity.Booking
		if err := db.Where("show_time_id = ? AND seat_id = ?", bookingRequest.ShowTimeID, seatID).First(&existingBooking).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Seat already booked for this showtime"})
			return
		}
	}

	// สร้างการจองใหม่สำหรับที่นั่งเหล่านั้น
	for _, seatID := range bookingRequest.SeatID {
		booking := entity.Booking{
			MemberID:    &bookingRequest.MemberID,
			ShowTimeID:  &bookingRequest.ShowTimeID,
			SeatID:      &seatID,
			BookingTime: time.Now().Format(time.RFC3339),
			Status:      "confirmed",
		}

		if err := db.Create(&booking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully"})
}
