package controller

import (
	//"fmt"
    "log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/tanapon395/sa-67-example/config"
	"github.com/tanapon395/sa-67-example/entity"
)

// GET /tickets
func ListTickets(c *gin.Context) {
	var tickets []entity.Ticket
	db := config.DB()

	// ดึงข้อมูลตั๋วทั้งหมด
	if err := db.Preload("ShowTime").Preload("Member").Find(&tickets).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

// GET /ticket/:id
func GetTicket(c *gin.Context) {
	ID := c.Param("id")
	var ticket entity.Ticket

	db := config.DB()
	// ดึงข้อมูลตั๋วที่มี ID ตรงกับที่ระบุ
	if err := db.Preload("ShowTime").Preload("Member").First(&ticket, ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

// ตัวอย่างการเขียน API ใน Golang เพื่อดึงข้อมูลจากหลายตาราง
func GetTicketsById(c *gin.Context) {
    memberID := c.Param("id")

    log.Println("Received Member ID:", memberID) // เพิ่มบรรทัดนี้เพื่อตรวจสอบค่า

    db := config.DB()
    var tickets []entity.Ticket
    if err := db.Preload("ShowTime.Movie").
        Preload("ShowTime.Theater").
        Preload("Seats").
        Where("member_id = ?", memberID).
        Order("created_at desc").
        Find(&tickets).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tickets not found"})
        return
    }

    log.Println("Fetched tickets:", tickets) // เพิ่มบรรทัดนี้เพื่อตรวจสอบค่า
    c.JSON(http.StatusOK, tickets)
}



// POST /tickets
func CreateTicket(c *gin.Context) {
	var ticket entity.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	// ตรวจสอบว่ามี ShowTime และ Member ที่ระบุไว้หรือไม่
	var showtime entity.ShowTimes
	if err := db.First(&showtime, ticket.ShowTimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}
	var member entity.Member
	if err := db.First(&member, ticket.MemberID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	// สร้างข้อมูลตั๋วใหม่
	if err := db.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created success", "data": ticket})
}

// PATCH /tickets/:id
func UpdateTicket(c *gin.Context) {
	var ticket entity.Ticket
	ID := c.Param("id")

	db := config.DB()
	// ค้นหาตั๋วตาม ID
	if err := db.First(&ticket, ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	// ทำการ bind ข้อมูลที่ได้รับจาก client ไปยังตัวแปร ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, unable to map payload"})
		return
	}

	// อัปเดตข้อมูลตั๋ว
	if err := db.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successful"})
}

// DELETE /tickets/:id
func DeleteTicket(c *gin.Context) {
	ID := c.Param("id")

	db := config.DB()
	if tx := db.Delete(&entity.Ticket{}, ID); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket ID not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successful"})
}
