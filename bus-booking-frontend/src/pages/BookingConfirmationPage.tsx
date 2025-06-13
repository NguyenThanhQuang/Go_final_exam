import React, { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import {
  Container,
  Card,
  Button,
  ListGroup,
  Alert,
  Spinner,
} from "react-bootstrap";
import type { Trip } from "../types/trip.types";
import { createBooking, type CreateBookingPayload } from "../api/bookingApi";
import { useAuth } from "../contexts/AuthContext";

interface LocationState {
  tripId: string;
  tripDetails: Trip;
  selectedSeatNumbers: string[];
  totalAmount: number;
}

const BookingConfirmationPage: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { isAuthenticated, isLoading: authLoading } = useAuth();

  const [bookingDetails, setBookingDetails] = useState<LocationState | null>(
    null
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  useEffect(() => {
    if (location.state) {
      const state = location.state as LocationState;
      if (
        state.tripId &&
        state.tripDetails &&
        state.selectedSeatNumbers &&
        state.totalAmount !== undefined
      ) {
        setBookingDetails(state);
      } else {
        setError("Dữ liệu đặt vé không hợp lệ. Vui lòng thử lại từ đầu.");
      }
    } else {
      setError("Không tìm thấy thông tin đặt vé. Vui lòng chọn lại chuyến đi.");
      setTimeout(() => navigate("/"), 3000);
    }
  }, [location.state, navigate]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate("/");
    }
  }, [isAuthenticated, authLoading, navigate]);

  const handleConfirmBooking = async () => {
    if (!bookingDetails || !isAuthenticated) {
      setError("Không thể xử lý yêu cầu. Vui lòng đăng nhập và thử lại.");
      return;
    }

    setLoading(true);
    setError(null);
    setSuccessMessage(null);

    const payload: CreateBookingPayload = {
      tripId: bookingDetails.tripId,
      seatNumbers: bookingDetails.selectedSeatNumbers,
    };

    try {
      const response = await createBooking(payload);
      if (response.dữ_liệu && response.dữ_liệu.id) {
        setSuccessMessage(
          response["thông báo"] ||
            "Giữ chỗ thành công! Chuẩn bị chuyển đến trang thanh toán."
        );
        setError(null);
        setTimeout(() => {
          navigate(`/payment/${response.dữ_liệu?.id}`, {
            state: {
              bookingId: response.dữ_liệu?.id,
              totalAmount: bookingDetails.totalAmount,
            },
          });
        }, 2000);
      } else {
        setError(
          response.lỗi ||
            response["thông báo"] ||
            "Giữ chỗ thất bại. Vui lòng thử lại."
        );
        setSuccessMessage(null);
      }
    } catch (err: any) {
      setError(
        err.lỗi || err["thông báo"] || "Có lỗi xảy ra trong quá trình giữ chỗ."
      );
      setSuccessMessage(null);
    } finally {
      setLoading(false);
    }
  };

  if (authLoading) {
    return (
      <Container className="text-center mt-5">
        <Spinner animation="border" />
      </Container>
    );
  }

  if (!bookingDetails && !error) {
    return (
      <Container className="text-center mt-5">
        <Spinner animation="border" />
      </Container>
    );
  }

  if (error && !bookingDetails) {
    return (
      <Container className="mt-5">
        <Alert variant="danger">
          {error}
          <hr />
          <Button onClick={() => navigate("/")} variant="primary">
            Quay về Trang chủ
          </Button>
        </Alert>
      </Container>
    );
  }

  const { tripDetails, selectedSeatNumbers, totalAmount } =
    bookingDetails || {};

  return (
    <Container className="mt-4" style={{ maxWidth: "800px" }}>
      <Card className="shadow">
        <Card.Header as="h3" className="text-center bg-primary text-white">
          Xác nhận Thông tin Đặt vé
        </Card.Header>
        <Card.Body>
          {error && (
            <Alert variant="danger" onClose={() => setError(null)} dismissible>
              {error}
            </Alert>
          )}
          {successMessage && <Alert variant="success">{successMessage}</Alert>}

          {tripDetails && (
            <>
              <h4 className="mb-3">Thông tin chuyến đi:</h4>
              <ListGroup variant="flush" className="mb-4">
                <ListGroup.Item>
                  <strong>Nhà xe:</strong> {tripDetails.companyName}
                </ListGroup.Item>
                <ListGroup.Item>
                  <strong>Tuyến:</strong> {tripDetails.route.from.name}{" "}
                  <i className="bi bi-arrow-right"></i>{" "}
                  {tripDetails.route.to.name}
                </ListGroup.Item>
                <ListGroup.Item>
                  <strong>Khởi hành:</strong>{" "}
                  {new Date(tripDetails.departureTime).toLocaleString()}
                </ListGroup.Item>
              </ListGroup>

              <h4 className="mb-3">Thông tin vé:</h4>
              <ListGroup variant="flush" className="mb-4">
                <ListGroup.Item>
                  <strong>Ghế đã chọn:</strong>{" "}
                  {selectedSeatNumbers?.join(", ")}
                </ListGroup.Item>
                <ListGroup.Item>
                  <strong>Số lượng:</strong> {selectedSeatNumbers?.length || 0}{" "}
                  vé
                </ListGroup.Item>
                <ListGroup.Item className="h5">
                  <strong>
                    Tổng thanh toán:{" "}
                    <span className="text-danger">
                      {(totalAmount || 0).toLocaleString()} VNĐ
                    </span>
                  </strong>
                </ListGroup.Item>
              </ListGroup>

              <div className="d-grid gap-2">
                <Button
                  variant="success"
                  size="lg"
                  onClick={handleConfirmBooking}
                  disabled={
                    loading ||
                    !!successMessage ||
                    !isAuthenticated ||
                    !bookingDetails
                  }
                >
                  {loading ? (
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  ) : (
                    "Xác nhận và Giữ chỗ"
                  )}
                </Button>
                <Button
                  variant="outline-secondary"
                  onClick={() => navigate(`/trips/${tripDetails.id}`)}
                  disabled={loading || !!successMessage || !bookingDetails}
                >
                  Chọn lại ghế
                </Button>
              </div>
            </>
          )}
        </Card.Body>
      </Card>
    </Container>
  );
};

export default BookingConfirmationPage;
