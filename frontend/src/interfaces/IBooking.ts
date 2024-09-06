export interface BookingInterface {
    ID?: number;              // ใช้เครื่องหมาย '?' เพื่อให้เป็น optional
    MemberID: number;          // รหัสสมาชิกที่ทำการจอง
    ShowTimeID: number;        // รหัสของรอบการฉาย
    SeatID: number[];          // Array ของรหัสที่นั่งที่ทำการจอง
    BookingTime?: Date;        // เวลาที่ทำการจอง (optional)
    Status: string;            // สถานะของการจอง (เช่น 'confirmed')
    Point?: number;            // จำนวนคะแนนที่ได้จากการจอง (optional)
  }
  