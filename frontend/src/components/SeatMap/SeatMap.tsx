import React, { useState, useEffect } from 'react';
import { Button, Select, message } from 'antd';
import axios from 'axios';
import Seat from './Seat';
import './SeatMap.css';
import { bookSeats } from '../../services/https'; // นำเข้าฟังก์ชันสำหรับการจอง

const { Option } = Select;

const seatsLeft = [
  ['A1', 'A2', 'A3', 'A4', 'A5', 'A6', 'A7', 'A8', 'A9', 'A10'],
  ['B1', 'B2', 'B3', 'B4', 'B5', 'B6', 'B7', 'B8', 'B9', 'B10'],
  ['C1', 'C2', 'C3', 'C4', 'C5', 'C6', 'C7', 'C8', 'C9', 'C10'],
  ['D1', 'D2', 'D3', 'D4', 'D5', 'D6', 'D7', 'D8', 'D9', 'D10'],
  ['E1', 'E2', 'E3', 'E4', 'E5', 'E6', 'E7', 'E8', 'E9', 'E10'],
  ['F1', 'F2', 'F3', 'F4', 'F5', 'F6', 'F7', 'F8', 'F9', 'F10'],
  ['G1', 'G2', 'G3', 'G4', 'G5', 'G6', 'G7', 'G8', 'G9', 'G10'],
  ['H1', 'H2', 'H3', 'H4', 'H5', 'H6', 'H7', 'H8', 'H9', 'H10'],
];

const seatsRight = [
  ['A11', 'A12', 'A13', 'A14', 'A15', 'A16', 'A17', 'A18', 'A19', 'A20'],
  ['B11', 'B12', 'B13', 'B14', 'B15', 'B16', 'B17', 'B18', 'B19', 'B20'],
  ['C11', 'C12', 'C13', 'C14', 'C15', 'C16', 'C17', 'C18', 'C19', 'C20'],
  ['D11', 'D12', 'D13', 'D14', 'D15', 'D16', 'D17', 'D18', 'D19', 'D20'],
  ['E11', 'E12', 'E13', 'E14', 'E15', 'E16', 'E17', 'E18', 'E19', 'E20'],
  ['F11', 'F12', 'F13', 'F14', 'F15', 'F16', 'F17', 'F18', 'F19', 'F20'],
  ['G11', 'G12', 'G13', 'G14', 'G15', 'G16', 'G17', 'G18', 'G19', 'G20'],
  ['H11', 'H12', 'H13', 'H14', 'H15', 'H16', 'H17', 'H18', 'H19', 'H20'],
];

const SeatMap: React.FC = () => {
  const [bookedSeats, setBookedSeats] = useState<string[]>([]);
  const [selectedSeats, setSelectedSeats] = useState<string[]>([]);
  const [memberID, setMemberID] = useState<number | null>(null);
  const [showtimeID, setShowtimeID] = useState<number>(1);
  const [theaterID, setTheaterID] = useState<number>(1);

  // ดึงข้อมูลจาก localStorage
  useEffect(() => {
    const storedMemberID = localStorage.getItem('memberID');
    const storedShowtimeID = localStorage.getItem('showtimeID');
    const storedTheaterID = localStorage.getItem('theaterID');

    if (storedMemberID) setMemberID(Number(storedMemberID));
    if (storedShowtimeID) setShowtimeID(Number(storedShowtimeID));
    if (storedTheaterID) setTheaterID(Number(storedTheaterID));
  }, []);

  // ดึงข้อมูลที่นั่งที่ถูกจองจาก API
  useEffect(() => {
    const fetchBookedSeats = async () => {
      try {
        const response = await fetch(`http://localhost:8000/api/booked-seats/${showtimeID}/${theaterID}`);
        const seatsData = await response.json();

        if (response.ok && seatsData.data) {
          const bookedSeatsArray = seatsData.data; // ใช้ข้อมูลที่ถูกส่งจาก API
          setBookedSeats(bookedSeatsArray); // อัปเดตสถานะที่นั่งที่ถูกจอง
        } else {
          console.error('Unexpected data structure:', seatsData);
        }
      } catch (error) {
        console.error('Error fetching booked seats:', error);
      }
    };

    fetchBookedSeats();
  }, [showtimeID, theaterID]); // ให้ทำงานเมื่อ showtimeID หรือ theaterID เปลี่ยนแปลง

  const onSelectSeat = (seat: string) => {
    setSelectedSeats(prev =>
      prev.includes(seat) ? prev.filter(s => s !== seat) : [...prev, seat]
    );
  };

  const handleConfirmBooking = async () => {
    if (!memberID) {
      message.error('Member ID not found');
      return;
    }

    const result = await bookSeats(showtimeID, theaterID, memberID, selectedSeats);

    if (result.success) {
      message.success(result.message);
      setSelectedSeats([]); // Reset selected seats
    } else {
      message.error(result.message);
    }
  };

  const totalPrice = selectedSeats.length * 100;

  return (
    <div className="container">
      <div style={{ textAlign: 'left', justifyContent: 'space-between', alignItems: 'flex-end', marginBottom: '20px', marginLeft: 200, marginRight: 200, display: 'flex' }}>
        <div>
          <h2 style={{ margin: 0 }}>SELECT TIME HERE!</h2>
          <div style={{ padding: 6 }}>
            <Select
              value={showtimeID} // เพิ่มค่า showtimeID ที่ดึงจาก localStorage
              style={{ width: 120, height: 40 }}
              onChange={(value) => setShowtimeID(value)} // เมื่อผู้ใช้เลือกเวลา จะอัปเดต showtimeID
            >
              <Option value={1}>12:00 AM</Option>
              <Option value={2}>14:00 PM</Option>
              <Option value={3}>02:40 PM</Option>
              <Option value={4}>04:00 PM</Option>
              <Option value={5}>06:20 PM</Option>
              <Option value={6}>08:40 PM</Option>
            </Select>
          </div>
          <div style={{ padding: 6 }}>
            <Select
              value={theaterID} // เพิ่มค่า theaterID ที่ดึงจาก localStorage
              style={{ width: 120, height: 40 }}
              onChange={(value) => setTheaterID(value)} // เมื่อผู้ใช้เลือกโรง จะอัปเดต theaterID
            >
              <Option value={1}>Theater 1</Option>
              <Option value={2}>Theater 2</Option>
              <Option value={3}>Theater 3</Option>
            </Select>
          </div>
        </div>
        <div style={{ textAlign: 'left', marginLeft: 'auto', alignItems: 'flex-end', display: 'flex' }}>
          <div style={{ display: 'flex', alignItems: 'flex-end', justifyContent: 'flex-end' }}>
            <div style={{ width: '15px', height: '15px', backgroundColor: '#1A2C50', marginRight: '5px', borderRadius: '3px' }}></div>
            <p style={{ margin: 0, marginRight: '10px' }}>Booked</p>
            <div style={{ width: '15px', height: '15px', backgroundColor: 'white', marginRight: '5px', borderRadius: '3px', border: '1px solid black' }}></div>
            <p style={{ margin: 0, marginRight: '10px' }}>Available</p>
            <div style={{ width: '15px', height: '15px', backgroundColor: '#007bff', marginRight: '5px', borderRadius: '3px' }}></div>
            <p style={{ margin: 0 }}>Selected</p>
          </div>
        </div>
      </div>
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <div>
          {seatsLeft.map((row, rowIndex) => (
            <div key={rowIndex} className={`seat-row ${rowIndex >= 6 ? 'yellow-border' : ''}`}>
              {row.map(seat => (
                <Seat
                  key={seat}
                  seat={seat}
                  isBooked={bookedSeats.includes(seat)}
                  isSelected={selectedSeats.includes(seat)}
                  onSelect={onSelectSeat}
                />
              ))}
            </div>
          ))}
        </div>
        <div style={{ width: '50px' }}></div> {/* ช่องว่างตรงกลาง */}
        <div>
          {seatsRight.map((row, rowIndex) => (
            <div key={rowIndex} className={`seat-row ${rowIndex >= 6 ? 'yellow-border' : ''}`}>
              {row.map(seat => (
                <Seat
                  key={seat}
                  seat={seat}
                  isBooked={bookedSeats.includes(seat)}
                  isSelected={selectedSeats.includes(seat)}
                  onSelect={onSelectSeat}
                />
              ))}
            </div>
          ))}
        </div>
      </div>
      <div className="beforesummary">
        <h2 style={{ color: 'black' }}>MERJE CINIPLEX</h2>
      </div>
      <div className="summary">
        <div className="summary-content">
          <h3>Total: Rp.Seat:</h3>
          <h3> {totalPrice}{selectedSeats.join(', ')}</h3>
          <div className="buttons">
            <Button type="default" style={{ marginRight: '10px' }}>BACK TO</Button>
            <Button type="primary" onClick={handleConfirmBooking}>CONFIRM</Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SeatMap;
