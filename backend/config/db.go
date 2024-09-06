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
        &entity.Booking{}, // เพิ่ม Booking ใน AutoMigrate
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

    // สร้างข้อมูลการฉายภาพยนตร์
    showTimes := []entity.ShowTimes{
        {Showdate: time.Date(2024, 10, 28, 14, 0, 0, 0, time.UTC), MovieID: &movies[0].ID, TheaterID: &theaters[0].ID},
        {Showdate: time.Date(2024, 10, 28, 16, 0, 0, 0, time.UTC), MovieID: &movies[0].ID, TheaterID: &theaters[1].ID},
        {Showdate: time.Date(2024, 10, 29, 12, 0, 0, 0, time.UTC), MovieID: &movies[2].ID, TheaterID: &theaters[2].ID},
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

    // สร้างข้อมูลที่นั่ง
    seats := []entity.Seat{
        {SeatNo: "A1", TheaterID: &theaters[0].ID},
        {SeatNo: "A2", TheaterID: &theaters[0].ID},
        {SeatNo: "A3", TheaterID: &theaters[0].ID},
    }

    for _, seat := range seats {
        if err := db.FirstOrCreate(&seat, entity.Seat{SeatNo: seat.SeatNo, TheaterID: seat.TheaterID}).Error; err != nil {
            fmt.Printf("Error creating Seat %s: %v\n", seat.SeatNo, err)
        }
    }

    // สร้างการจองใหม่
    seatNumbersForBooking := []string{"A1", "A2", "A3"}
    for _, seatNo := range seatNumbersForBooking {
        var seat entity.Seat
        // ตรวจสอบว่าที่นั่งมีอยู่ในฐานข้อมูลหรือไม่
        if err := db.Where("seat_no = ? AND theater_id = ?", seatNo, theaters[0].ID).First(&seat).Error; err != nil {
            fmt.Printf("Seat not found for SeatNo: %s, TheaterID: %d. Error: %v\n", seatNo, theaters[0].ID, err)
            continue
        }

        // สร้างการจอง
        booking := entity.Booking{
            MemberID:    &Member2.ID,
            ShowTimeID:  &showTimes[0].ID,
            SeatID:      &seat.ID,
            BookingTime: time.Now().String(),
            Status:      "confirmed",
        }

        if err := db.Create(&booking).Error; err != nil {
            fmt.Printf("Error creating booking for Seat %s: %v\n", seatNo, err)
            continue
        }

        // แสดงการจองที่สร้างสำเร็จ
        fmt.Printf("Booking created for Seat %s in ShowTime %d with Booking ID %d\n", seatNo, showTimes[0].ID, booking.ID)

        // อัปเดต Ticket ให้เชื่อมโยงกับ BookingID
        if err := db.Model(&entity.Ticket{}).Where("member_id = ? AND show_time_id = ?", Member2.ID, showTimes[0].ID).Update("booking_id", booking.ID).Error; err != nil {
            fmt.Printf("Error updating Ticket with Booking ID for Seat %s: %v\n", seatNo, err)
            continue
        }

        fmt.Printf("Ticket updated for Seat %s with Booking ID %d\n", seatNo, booking.ID)
    }


	// เพิ่มการจองสำหรับรอบที่สอง (จองที่นั่ง A1, A2)
seatNumbersForBookingRound2 := []string{"A1", "A2"}
for _, seatNo := range seatNumbersForBookingRound2 {
    var seat entity.Seat
    // ตรวจสอบว่าที่นั่งมีอยู่ในฐานข้อมูลหรือไม่
    if err := db.Where("seat_no = ? AND theater_id = ?", seatNo, theaters[0].ID).First(&seat).Error; err != nil {
        fmt.Printf("Seat not found for SeatNo: %s, TheaterID: %d. Error: %v\n", seatNo, theaters[0].ID, err)
        continue
    }

    // สร้างการจองสำหรับรอบที่ 2
    bookingRound2 := entity.Booking{
        MemberID:    &Member2.ID,
        ShowTimeID:  &showTimes[1].ID,  // ใช้ ShowTimeID ของรอบที่ 2
        SeatID:      &seat.ID,
        BookingTime: time.Now().String(),
        Status:      "confirmed",
    }

    if err := db.Create(&bookingRound2).Error; err != nil {
        fmt.Printf("Error creating booking for Seat %s in round 2: %v\n", seatNo, err)
        continue
    }

    // แสดงการจองที่สร้างสำเร็จ
    fmt.Printf("Booking created for Seat %s in ShowTime %d with Booking ID %d in round 2\n", seatNo, showTimes[1].ID, bookingRound2.ID)

    // อัปเดต Ticket ให้เชื่อมโยงกับ BookingID ของรอบที่ 2
    if err := db.Model(&entity.Ticket{}).Where("member_id = ? AND show_time_id = ?", Member2.ID, showTimes[1].ID).Update("booking_id", bookingRound2.ID).Error; err != nil {
        fmt.Printf("Error updating Ticket with Booking ID for Seat %s in round 2: %v\n", seatNo, err)
        continue
    }

    fmt.Printf("Ticket updated for Seat %s with Booking ID %d in round 2\n", seatNo, bookingRound2.ID)
}


    fmt.Println("Database setup complete")
}
