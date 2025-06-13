import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Container,
  Row,
  Col,
  Card,
  Button,
  Spinner,
  Alert,
  ListGroup,
} from "react-bootstrap";
import { getTripDetails } from "../api/tripApi";
import type { Trip, Seat } from "../types/trip.types";
import SeatButton from "../components/trips/SeatButton";
import { useAuth } from "../contexts/AuthContext";

const TripDetailsPage: React.FC = () => {
  const { tripId } = useParams<{ tripId: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading } = useAuth();

  const [trip, setTrip] = useState<Trip | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedSeats, setSelectedSeats] = useState<string[]>([]);

  useEffect(() => {
    if (tripId) {
      const fetchTripDetails = async () => {
        setLoading(true);
        setError(null);
        try {
          const data = await getTripDetails(tripId);
          if (data) {
            setTrip(data);
          } else {
            setError("Không tìm thấy thông tin chuyến đi.");
          }
        } catch (err) {
          setError("Lỗi tải dữ liệu chuyến đi.");
        } finally {
          setLoading(false);
        }
      };
      fetchTripDetails();
    } else {
      setError("ID chuyến đi không hợp lệ.");
      setLoading(false);
    }
  }, [tripId]);

  const handleSeatClick = (seatNumber: string) => {
    setSelectedSeats((prevSelected) =>
      prevSelected.includes(seatNumber)
        ? prevSelected.filter((s) => s !== seatNumber)
        : [...prevSelected, seatNumber]
    );
  };

  const calculateTotalPrice = () => {
    if (!trip) return 0;
    return selectedSeats.length * trip.price;
  };

  const handleProceedToBooking = () => {
    if (!isAuthenticated) {
      alert("Vui lòng đăng nhập để tiếp tục đặt vé.");
      return;
    }
    if (selectedSeats.length === 0) {
      alert("Vui lòng chọn ít nhất một ghế.");
      return;
    }
    navigate("/booking/confirm", {
      state: {
        tripId: trip?.id,
        tripDetails: trip,
        selectedSeatNumbers: selectedSeats,
        totalAmount: calculateTotalPrice(),
      },
    });
  };

  if (authLoading || loading) {
    return (
      <Container className="text-center mt-5">
        <Spinner animation="border" variant="primary" />
        <p>Đang tải dữ liệu...</p>
      </Container>
    );
  }

  if (error) {
    return (
      <Container className="mt-5">
        <Alert variant="danger">{error}</Alert>
        <Button onClick={() => navigate("/")}>Quay về Trang chủ</Button>
      </Container>
    );
  }

  if (!trip) {
    return (
      <Container className="mt-5">
        <Alert variant="info">
          Không có thông tin chi tiết cho chuyến đi này.
        </Alert>
      </Container>
    );
  }

  const seatsPerRow = 4;
  const seatRows: Seat[][] = [];
  if (trip.seats) {
    for (let i = 0; i < trip.seats.length; i += seatsPerRow) {
      seatRows.push(trip.seats.slice(i, i + seatsPerRow));
    }
  }

  return (
    <Container className="mt-4">
      <Row>
        <Col md={8}>
          <Card className="mb-4 shadow-sm">
            <Card.Header as="h4" className="bg-primary text-white">
              Chi tiết chuyến đi: {trip.companyName}
            </Card.Header>
            <Card.Body>
              <Card.Title>
                {trip.route.from.name}{" "}
                <i className="bi bi-arrow-right-circle-fill"></i>{" "}
                {trip.route.to.name}
              </Card.Title>
              <ListGroup variant="flush">
                <ListGroup.Item>
                  <strong>Khởi hành:</strong>{" "}
                  {new Date(trip.departureTime).toLocaleString()}
                </ListGroup.Item>
                <ListGroup.Item>
                  <strong>Dự kiến đến:</strong>{" "}
                  {new Date(trip.expectedArrivalTime).toLocaleString()}
                </ListGroup.Item>
                <ListGroup.Item>
                  <strong>Giá vé:</strong> {trip.price.toLocaleString()} VNĐ/ghế
                </ListGroup.Item>
                <ListGroup.Item>
                  <strong>Số ghế trống:</strong> {trip.availableSeats}
                </ListGroup.Item>
              </ListGroup>
            </Card.Body>
          </Card>

          <Card className="shadow-sm">
            <Card.Header as="h5" className="bg-light">
              Chọn ghế của bạn
            </Card.Header>
            <Card.Body className="text-center">
              <p className="text-muted">
                <span className="me-3">
                  <Button
                    size="sm"
                    variant="outline-success"
                    disabled
                    className="p-1"
                  ></Button>{" "}
                  Trống
                </span>
                <span className="me-3">
                  <Button
                    size="sm"
                    variant="primary"
                    disabled
                    className="p-1"
                  ></Button>{" "}
                  Đang chọn
                </span>
                <span className="me-3">
                  <Button
                    size="sm"
                    variant="secondary"
                    disabled
                    className="p-1"
                  ></Button>{" "}
                  Đã đặt
                </span>
                <span>
                  <Button
                    size="sm"
                    variant="warning"
                    disabled
                    className="p-1"
                  ></Button>{" "}
                  Đang giữ
                </span>
              </p>
              <hr />
              {seatRows.map((row, rowIndex) => (
                <div
                  key={rowIndex}
                  className="mb-2 d-flex justify-content-center"
                >
                  {row.map((seat) => (
                    <SeatButton
                      key={seat.seatNumber}
                      seat={seat}
                      isSelected={selectedSeats.includes(seat.seatNumber)}
                      onSeatClick={handleSeatClick}
                    />
                  ))}
                </div>
              ))}
              {!trip.seats ||
                (trip.seats.length === 0 && (
                  <p>Không có thông tin sơ đồ ghế.</p>
                ))}
            </Card.Body>
          </Card>
        </Col>

        <Col md={4}>
          <Card className="sticky-top shadow-sm" style={{ top: "80px" }}>
            <Card.Header as="h5" className="bg-success text-white">
              Tóm tắt đặt vé
            </Card.Header>
            <Card.Body>
              {selectedSeats.length > 0 ? (
                <>
                  <p>
                    <strong>Ghế đã chọn:</strong> {selectedSeats.join(", ")}
                  </p>
                  <p>
                    <strong>Số lượng:</strong> {selectedSeats.length} ghế
                  </p>
                  <hr />
                  <h5>
                    <strong>Tổng tiền:</strong>{" "}
                    {calculateTotalPrice().toLocaleString()} VNĐ
                  </h5>
                </>
              ) : (
                <p>Vui lòng chọn ghế từ sơ đồ bên trái.</p>
              )}
              <Button
                variant="primary"
                className="w-100 mt-3"
                onClick={handleProceedToBooking}
                disabled={selectedSeats.length === 0 || !isAuthenticated}
              >
                {isAuthenticated ? "Tiến hành đặt vé" : "Đăng nhập để đặt vé"}
              </Button>
              {!isAuthenticated && (
                <small className="d-block text-center mt-2 text-danger">
                  Bạn cần đăng nhập để thực hiện chức năng này.
                </small>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </Container>
  );
};

export default TripDetailsPage;
