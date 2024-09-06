package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanapon395/sa-67-example/config"
	"github.com/tanapon395/sa-67-example/entity"
)

// GET /seats
func ListSeats(c *gin.Context) {
	var seats []entity.Seat
	db := config.DB()

	// ดึงข้อมูลที่นั่งทั้งหมด
	if err := db.Preload("Theater").Find(&seats).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seats)
}

// GET /seat/:id
func GetSeat(c *gin.Context) {
	ID := c.Param("id")
	var seat entity.Seat

	db := config.DB()
	// ดึงข้อมูลที่นั่งที่มี ID ตรงกับที่ระบุ
	if err := db.Preload("Theater").First(&seat, ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seat)
}

// POST /seats
func CreateSeat(c *gin.Context) {
	var seat entity.Seat
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	// ตรวจสอบว่ามีโรงหนังที่ระบุไว้หรือไม่
	var theater entity.Theater
	if err := db.First(&theater, seat.TheaterID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "theater not found"})
		return
	}

	// สร้างข้อมูลที่นั่งใหม่
	if err := db.Create(&seat).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created success", "data": seat})
}

// PATCH /seats/:id
func UpdateSeat(c *gin.Context) {
	var seat entity.Seat
	ID := c.Param("id")

	db := config.DB()
	// ค้นหาที่นั่งตาม ID
	if err := db.First(&seat, ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "seat not found"})
		return
	}

	// ทำการ bind ข้อมูลที่ได้รับจาก client ไปยังตัวแปร seat
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, unable to map payload"})
		return
	}

	// อัปเดตข้อมูลที่นั่ง
	if err := db.Save(&seat).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successful"})
}

// DELETE /seats/:id
func DeleteSeat(c *gin.Context) {
	ID := c.Param("id")

	db := config.DB()
	if tx := db.Delete(&entity.Seat{}, ID); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seat ID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successful"})
}

// GET /booked-seats/:showtimeID
// GET /booked-seats/:showtimeID/:theaterID
// GET /booked-seats/:showtimeID/:theaterID
// GET /booked-seats/:showtimeID/:theaterID
// GET /booked-seats/:showtimeID/:theaterID
func GetBookedSeats(c *gin.Context) {
    var bookings []entity.Booking
    showtimeID := c.Param("showtimeID")

    // ดึงข้อมูลการจองที่ตรงกับ showtimeID
    if err := config.DB().
        Where("show_time_id = ?", showtimeID).
        Find(&bookings).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No booked seats found"})
        return
    }

    // สร้างรายการที่นั่งที่ถูกจองแล้วจากตาราง bookings
    bookedSeats := make([]string, 0)
    for _, booking := range bookings {
        var seat entity.Seat
        if err := config.DB().First(&seat, booking.SeatID).Error; err == nil {
            bookedSeats = append(bookedSeats, seat.SeatNo)
        }
    }

    c.JSON(http.StatusOK, gin.H{"data": bookedSeats})
}




// POST /book-seats
func BookSeats(c *gin.Context) {
    var bookingData struct {
        Seats      []string `json:"seats"`
        ShowtimeID uint     `json:"showtimeID"`
        TheaterID  uint     `json:"theaterID"`
        MemberID   uint     `json:"memberID"`
    }

    if err := c.ShouldBindJSON(&bookingData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    db := config.DB()

    // Create a new Ticket entry
    ticket := entity.Ticket{
        ShowTimeID: &bookingData.ShowtimeID,
        MemberID:   &bookingData.MemberID,
        // Assume points calculation and payment information is handled elsewhere
    }

    if err := db.Create(&ticket).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
        return
    }

    // Update the seats with the new TicketID
    for _, seatNo := range bookingData.Seats {
        if err := db.Model(&entity.Seat{}).
            Where("seat_no = ? AND theater_id = ?", seatNo, bookingData.TheaterID).
            Update("ticket_id", ticket.ID).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book seats"})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"message": "Seats booked successfully"})
}

