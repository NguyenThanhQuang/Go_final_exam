import React, { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import {
  Container,
  Card,
  Spinner,
  Alert,
  ListGroup,
  Button,
  Badge,
  Row,
  Col,
} from "react-bootstrap";
import { getMyBookings } from "../api/bookingApi";
import type { Booking } from "../types/booking.types";
import { useAuth } from "../contexts/AuthContext";

const BookingHistoryPage: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading } = useAuth();

  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate("/?showLogin=true");
      return;
    }

    if (isAuthenticated) {
      const fetchBookings = async () => {
        setLoading(true);
        setError(null);
        try {
          const data = await getMyBookings();
          setBookings(data);
        } catch (err: any) {
          setError(
            err.lỗi || err["thông báo"] || "Lỗi khi tải lịch sử đặt vé."
          );
        } finally {
          setLoading(false);
        }
      };
      fetchBookings();
    }
  }, [isAuthenticated, authLoading, navigate]);

  const getStatusBadgeVariant = (status: string = "pending") => {
    switch (status.toLowerCase()) {
      case "confirmed":
        return "success";
      case "held":
        return "warning";
      case "cancelled":
      case "expired":
        return "danger";
      default:
        return "secondary";
    }
  };

  if (authLoading || loading) {
    return (
      <Container className="text-center mt-5">
        <Spinner animation="border" variant="info" />
        <p>Đang tải lịch sử vé của bạn...</p>
      </Container>
    );
  }

  if (error) {
    return (
      <Container className="mt-5">
        <Alert variant="danger" className="text-center">
          <h4>Lỗi!</h4>
          <p>{error}</p>
          <Link to="/" className="btn btn-primary">
            Quay về Trang chủ
          </Link>
        </Alert>
      </Container>
    );
  }

  return (
    <Container className="my-4">
      <h2 className="mb-4 text-center">Lịch sử Đặt vé của bạn</h2>
      {bookings.length === 0 ? (
        <Alert variant="info" className="text-center">
          Bạn chưa có giao dịch đặt vé nào.
          <hr />
          <Link to="/" className="btn btn-primary">
            Tìm chuyến đi ngay!
          </Link>
        </Alert>
      ) : (
        <ListGroup>
          {bookings.map((booking) => (
            <ListGroup.Item
              key={booking.id}
              className="mb-3 shadow-sm booking-history-item"
            >
              <Row className="align-items-center">
                <Col md={8}>
                  <h5>
                    {booking.tripInfo?.route.from.name || "N/A"}{" "}
                    <i className="bi bi-arrow-right"></i>{" "}
                    {booking.tripInfo?.route.to.name || "N/A"}
                  </h5>
                  <p className="mb-1">
                    <strong>Nhà xe:</strong>{" "}
                    {booking.tripInfo?.companyName || "N/A"}
                  </p>
                  <p className="mb-1">
                    <strong>Ngày đi:</strong>{" "}
                    {booking.tripInfo
                      ? new Date(
                          booking.tripInfo.departureTime
                        ).toLocaleDateString()
                      : "N/A"}
                    {" - "}
                    <strong>Giờ đi:</strong>{" "}
                    {booking.tripInfo
                      ? new Date(
                          booking.tripInfo.departureTime
                        ).toLocaleTimeString([], {
                          hour: "2-digit",
                          minute: "2-digit",
                        })
                      : "N/A"}
                  </p>
                  <p className="mb-1">
                    <strong>Ghế:</strong>{" "}
                    {booking.passengers?.map((p) => p.seatNumber).join(", ") ||
                      "N/A"}
                  </p>
                  <p className="mb-0">
                    <strong>Ngày đặt:</strong>{" "}
                    {new Date(booking.bookingTime).toLocaleString()}
                  </p>
                </Col>
                <Col md={4} className="text-md-end mt-3 mt-md-0">
                  <p className="mb-2">
                    <Badge
                      bg={getStatusBadgeVariant(booking.status)}
                      pill
                      className="fs-6 px-3 py-2"
                    >
                      {booking.status.toUpperCase()}
                    </Badge>
                  </p>
                  <h5 className="text-danger mb-2">
                    {booking.totalAmount.toLocaleString()} VNĐ
                  </h5>
                  <Link
                    to={`/ticket/${booking.id}`}
                    className="btn btn-info btn-sm"
                  >
                    Xem chi tiết vé
                  </Link>
                </Col>
              </Row>
            </ListGroup.Item>
          ))}
        </ListGroup>
      )}
      <style jsx global>{`
        .booking-history-item {
          border-radius: 0.375rem;
          border: 1px solid #dee2e6;
          transition: box-shadow 0.15s ease-in-out;
        }
        .booking-history-item:hover {
          box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.1) !important;
        }
        .fs-6 {
          font-size: 0.9rem !important;
        }
      `}</style>
    </Container>
  );
};

export default BookingHistoryPage;
