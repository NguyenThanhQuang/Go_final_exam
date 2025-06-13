import React, { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import {
  Container,
  Card,
  Spinner,
  Alert,
  ListGroup,
  Row,
  Col,
  Badge,
  Button,
} from "react-bootstrap";
import { getBookingDetails } from "../api/bookingApi";
import type { Booking } from "../types/booking.types";
import { useAuth } from "../contexts/AuthContext";

const TicketDetailsPage: React.FC = () => {
  const { bookingId } = useParams<{ bookingId: string }>();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading } = useAuth();

  const [booking, setBooking] = useState<Booking | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      console.log(
        "TicketDetailsPage: Chưa đăng nhập hoặc auth chưa load, navigating to /"
      );
      navigate("/?showLogin=true"); // Redirect và có thể yêu cầu hiện modal login
      return;
    }

    if (bookingId && isAuthenticated) {
      console.log(
        `TicketDetailsPage: Fetching details for bookingId: ${bookingId}`
      );
      const fetchBookingDetails = async () => {
        setLoading(true);
        setError(null);
        try {
          const data = await getBookingDetails(bookingId);
          console.log(
            "TicketDetailsPage: Dữ liệu booking nhận được từ API:",
            data
          );

          if (data) {
            setBooking(data);
            if (data.tripInfo) {
              console.log(
                "TicketDetailsPage: tripInfo tồn tại:",
                data.tripInfo
              );
              if (data.tripInfo.route) {
                console.log(
                  "TicketDetailsPage: tripInfo.route tồn tại:",
                  data.tripInfo.route
                );
              } else {
                console.warn(
                  "TicketDetailsPage: tripInfo.route KHÔNG tồn tại trong tripInfo."
                );
              }
            } else {
              console.warn(
                "TicketDetailsPage: tripInfo KHÔNG tồn tại hoặc undefined trong booking data."
              );
            }
          } else {
            setError(
              `Không tìm thấy thông tin cho vé #${bookingId}. Vé có thể không tồn tại hoặc bạn không có quyền xem.`
            );
          }
        } catch (err: any) {
          console.error(
            "TicketDetailsPage: Lỗi khi gọi getBookingDetails:",
            err
          );
          setError(err.lỗi || err["thông báo"] || "Lỗi khi tải thông tin vé.");
        } finally {
          setLoading(false);
        }
      };
      fetchBookingDetails();
    } else if (!bookingId) {
      setError("ID vé không hợp lệ.");
      setLoading(false);
    }
    // Không fetch nếu !isAuthenticated, vì useEffect trên sẽ xử lý redirect
  }, [bookingId, isAuthenticated, authLoading, navigate]);

  if (authLoading || loading) {
    return (
      <Container className="text-center mt-5">
        <Spinner animation="border" variant="success" />
        <p className="mt-2">Đang tải thông tin vé của bạn...</p>
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

  if (!booking) {
    return (
      <Container className="mt-5">
        <Alert variant="warning" className="text-center">
          Không có thông tin chi tiết cho vé này.
          <hr />
          <Link to="/my-bookings" className="btn btn-info me-2">
            Xem vé của tôi
          </Link>
          <Link to="/" className="btn btn-outline-secondary">
            Quay về Trang chủ
          </Link>
        </Alert>
      </Container>
    );
  }

  // Lấy thông tin chuyến đi từ booking.tripInfo (đã được backend populate)
  const trip = booking.tripInfo;
  console.log(
    "TicketDetailsPage: Giá trị của 'trip' (booking.tripInfo) khi render:",
    trip
  );

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

  return (
    <Container className="my-4">
      <Card className="shadow-lg ticket-card">
        <Card.Header className="bg-success text-white p-3">
          <Row className="align-items-center">
            <Col>
              <h2 className="mb-0">Vé Xe Điện Tử</h2>
            </Col>
            <Col xs="auto" className="text-end">
              Mã vé:{" "}
              <strong className="fs-5">
                #{booking.ticketCode || booking.id.slice(-8).toUpperCase()}
              </strong>
            </Col>
          </Row>
        </Card.Header>
        <Card.Body className="p-4">
          <Row className="mb-4">
            <Col md={7}>
              <h4 className="text-primary">
                {trip?.companyName || "N/A (Không có tên nhà xe)"}
              </h4>
              <h5>
                {trip?.route?.from?.name || "N/A (Điểm đi)"}{" "}
                <i className="bi bi-arrow-long-right"></i>{" "}
                {trip?.route?.to?.name || "N/A (Điểm đến)"}
              </h5>
              <p className="mb-1">
                <strong>Ngày đi:</strong>{" "}
                {trip?.departureTime
                  ? new Date(trip.departureTime).toLocaleDateString("vi-VN")
                  : "N/A"}
              </p>
              <p className="mb-1">
                <strong>Giờ khởi hành:</strong>{" "}
                {trip?.departureTime
                  ? new Date(trip.departureTime).toLocaleTimeString("vi-VN", {
                      hour: "2-digit",
                      minute: "2-digit",
                    })
                  : "N/A"}
              </p>
              <p>
                <strong>Dự kiến đến:</strong>{" "}
                {trip?.expectedArrivalTime
                  ? new Date(trip.expectedArrivalTime).toLocaleTimeString(
                      "vi-VN",
                      { hour: "2-digit", minute: "2-digit" }
                    )
                  : "N/A"}{" "}
                (
                {trip?.expectedArrivalTime
                  ? new Date(trip.expectedArrivalTime).toLocaleDateString(
                      "vi-VN"
                    )
                  : ""}
                )
              </p>
            </Col>
            <Col md={5} className="text-md-end mt-3 mt-md-0">
              <p className="mb-1">
                <strong>Trạng thái vé:</strong>{" "}
                <Badge
                  bg={getStatusBadgeVariant(booking.status)}
                  pill
                  className="fs-6 px-3 py-2"
                >
                  {booking.status ? booking.status.toUpperCase() : "N/A"}
                </Badge>
              </p>
              <p className="mb-1">
                <strong>Ngày đặt:</strong>{" "}
                {new Date(booking.bookingTime).toLocaleString("vi-VN")}
              </p>
              {/* Giả định paymentStatus luôn là "paid" vì chúng ta giả lập thành công ở frontend */}
              {/* Trong thực tế, bạn sẽ dùng booking.paymentStatus từ backend */}
              <p>
                <strong>Thanh toán:</strong>{" "}
                <span className="text-success fw-bold">
                  Đã thanh toán (Giả lập)
                </span>
              </p>
            </Col>
          </Row>

          <hr />

          <h5 className="mt-4 mb-3">Thông tin hành khách và ghế:</h5>
          <ListGroup>
            {booking.passengers && booking.passengers.length > 0 ? (
              booking.passengers.map((passenger, index) => (
                <ListGroup.Item
                  key={index}
                  className="d-flex justify-content-between align-items-center"
                >
                  <div>
                    <strong>Hành khách {index + 1}:</strong>{" "}
                    {passenger.name || "Chưa cập nhật tên"}
                    {passenger.phone && ` - SĐT: ${passenger.phone}`}
                  </div>
                  <span className="fw-bold">Ghế: {passenger.seatNumber}</span>
                </ListGroup.Item>
              ))
            ) : (
              <ListGroup.Item>
                Không có thông tin hành khách chi tiết.
              </ListGroup.Item>
            )}
          </ListGroup>

          <hr className="my-4" />

          <Row>
            <Col>
              <h4 className="text-danger">
                Tổng tiền:{" "}
                {booking.totalAmount
                  ? booking.totalAmount.toLocaleString("vi-VN")
                  : "0"}{" "}
                VNĐ
              </h4>
            </Col>
          </Row>

          <Alert variant="info" className="mt-4">
            <strong>Lưu ý:</strong> Vui lòng có mặt tại điểm đón trước giờ khởi
            hành ít nhất 30 phút. Mang theo thông tin vé này (bản điện tử hoặc
            bản in) để đối chiếu khi lên xe.
          </Alert>
        </Card.Body>
        <Card.Footer className="text-center p-3">
          <Link to="/my-bookings" className="btn btn-outline-primary me-2">
            Xem tất cả vé
          </Link>
          <Link to="/" className="btn btn-outline-secondary">
            Tìm chuyến đi khác
          </Link>
        </Card.Footer>
      </Card>
      <style>{`
        .ticket-card {
          border: 1px solid #198754; 
          border-radius: 0.5rem;
        }
        .bi-arrow-long-right {
          font-size: 1.2em;
          vertical-align: middle;
        }
        .fs-6 {
          font-size: 0.9rem !important;
        }
      `}</style>
    </Container>
  );
};

export default TicketDetailsPage;
