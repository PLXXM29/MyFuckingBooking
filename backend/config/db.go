package config

import (
	"fmt"
	"time"

	"github.com/tanapon395/sa-67-example/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}

func ConnectionDB() {
	database, err := gorm.Open(sqlite.Open("sa.db?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("connected database")
	db = database
}

func SetupDatabase() {

	db.AutoMigrate(
		&entity.Member{},
		&entity.Gender{},
		&entity.Movie{},
		&entity.Theater{},
		&entity.ShowTimes{},
		&entity.Ticket{},
		&entity.Seat{},
		&entity.Payment{},
	)

	// สร้างข้อมูลเพศ
	GenderMale := entity.Gender{Name: "Male"}
	GenderFemale := entity.Gender{Name: "Female"}

	db.FirstOrCreate(&GenderMale, &entity.Gender{Name: "Male"})
	db.FirstOrCreate(&GenderFemale, &entity.Gender{Name: "Female"})

	// สร้างข้อมูลสมาชิก
	hashedPassword, _ := HashPassword("123456")
	Member := &entity.Member{
		FirstName:  "Software",
		LastName:   "Analysis",
		Email:      "sa@gmail.com",
		Password:   hashedPassword,
		GenderID:   GenderMale.ID,
		TotalPoint: 2,
	}
	db.FirstOrCreate(Member, &entity.Member{
		Email: "sa@gmail.com",
	})

	hashedPassword2, _ := HashPassword("123456")
	Member2 := &entity.Member{
		FirstName:  "Software2",
		LastName:   "Analysis2",
		Email:      "sa2@gmail.com",
		Password:   hashedPassword2,
		GenderID:   GenderMale.ID,
		TotalPoint: 5,
	}

	db.FirstOrCreate(Member2, &entity.Member{
		Email: "sa2@gmail.com",
	})

	// สร้างข้อมูลภาพยนตร์ 3 เรื่อง
	movies := []entity.Movie{
		{MovieName: "Inception", MovieDuration: 70},
		{MovieName: "The Dark Knight", MovieDuration: 80},
		{MovieName: "Interstellar", MovieDuration: 90},
	}

	for i := range movies {
		db.FirstOrCreate(&movies[i], entity.Movie{MovieName: movies[i].MovieName})
	}

	// สร้างข้อมูลโรงหนัง 3 โรง
	theaters := []entity.Theater{
		{TheaterName: "Theater 1"},
		{TheaterName: "Theater 2"},
		{TheaterName: "Theater 3"},
	}

	for i := range theaters {
		db.FirstOrCreate(&theaters[i], entity.Theater{TheaterName: theaters[i].TheaterName})
	}

	// สร้างที่นั่งสำหรับแต่ละโรงหนัง
	seatNumbers := []string{}
	for row := 'A'; row <= 'H'; row++ {
		for num := 1; num <= 20; num++ {
			seatNumber := fmt.Sprintf("%c%d", row, num)
			seatNumbers = append(seatNumbers, seatNumber)
		}
	}

	for _, theater := range theaters {
		for _, seatNo := range seatNumbers {
			seat := entity.Seat{
				SeatNo:    seatNo,
				Status:    "Available",
				Price:     200,
				TheaterID: &theater.ID,
			}
			db.FirstOrCreate(&seat, &entity.Seat{SeatNo: seatNo, TheaterID: &theater.ID})
		}
	}

	// สร้างข้อมูล show_times
	showTimes := []entity.ShowTimes{
		{Showdate: time.Date(2024, 10, 28, 12, 0, 0, 0, time.UTC), MovieID: &movies[0].ID, TheaterID: &theaters[0].ID}, // movie_id = 1, theater_id = 1
		{Showdate: time.Date(2024, 10, 28, 14, 0, 0, 0, time.UTC), MovieID: &movies[0].ID, TheaterID: &theaters[0].ID}, // movie_id = 1, theater_id = 1
		{Showdate: time.Date(2024, 10, 29, 12, 0, 0, 0, time.UTC), MovieID: &movies[2].ID, TheaterID: &theaters[2].ID}, // movie_id = 3, theater_id = 3
	}

	for i := range showTimes {
		db.FirstOrCreate(&showTimes[i], entity.ShowTimes{Showdate: showTimes[i].Showdate, MovieID: showTimes[i].MovieID, TheaterID: showTimes[i].TheaterID})
	}

	// สร้าง tickets สำหรับสมาชิกที่ 2
	tickets := []entity.Ticket{
		{Point: 10, ShowTimeID: &showTimes[0].ID, MemberID: &Member2.ID},
		{Point: 15, ShowTimeID: &showTimes[1].ID, MemberID: &Member2.ID},
	}

	for i := range tickets {
		db.FirstOrCreate(&tickets[i], entity.Ticket{ShowTimeID: tickets[i].ShowTimeID, MemberID: tickets[i].MemberID})
	}

	// สร้าง payments และเชื่อมโยง ticket_id ที่ถูกต้อง
	payments := []entity.Payment{
		{TotalPrice: 600, Status: "Paid", PaymentTime: time.Now(), MemberID: &Member2.ID, TicketID: &tickets[0].ID},
		{TotalPrice: 900, Status: "Paid", PaymentTime: time.Now(), MemberID: &Member2.ID, TicketID: &tickets[1].ID},
	}

	for i := range payments {
		db.Create(&payments[i])
		fmt.Printf("Payment %d created with ID: %d\n", i+1, payments[i].ID)

		// อัปเดต ticket ให้เชื่อมโยงกับ payment ที่ถูกต้อง
		db.Model(&tickets[i]).Update("PaymentID", payments[i].ID)
		if db.Error != nil {
			fmt.Printf("Error updating Ticket %d with Payment ID: %v\n", tickets[i].ID, db.Error)
		}
	}

	// อัปเดตที่นั่งที่ถูกจองในโรงหนัง
	seatNumbersForBooking := []string{"A1", "A2", "A3"}
	for i, ticket := range tickets {
		for _, seatNo := range seatNumbersForBooking {
			db.Model(&entity.Seat{}).Where("seat_no = ? AND theater_id = ?", seatNo, showTimes[i].TheaterID).Update("ticket_id", ticket.ID)
		}
	}

	fmt.Println("Database setup complete")
}
