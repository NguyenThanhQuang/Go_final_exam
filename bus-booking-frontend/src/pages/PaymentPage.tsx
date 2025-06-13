import React, { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import { Container, Card, Spinner, Alert, Button } from "react-bootstrap";
import { useAuth } from "../contexts/AuthContext";

const PaymentPage: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { bookingId: bookingIdFromParams } = useParams<{ bookingId: string }>();
  const { isAuthenticated, isLoading: authLoading } = useAuth();

  const [paymentStatus, setPaymentStatus] = useState<
    "pending" | "processing" | "success" | "failed"
  >("pending");
  const [error, setError] = useState<string | null>(null);
  const [bookingInfo, setBookingInfo] = useState<{
    bookingId: string;
    totalAmount: number;
  } | null>(null);

  useEffect(() => {
    const currentBookingId = bookingIdFromParams || location.state?.bookingId;
    const currentTotalAmount = location.state?.totalAmount;

    if (currentBookingId && currentTotalAmount !== undefined) {
      setBookingInfo({
        bookingId: currentBookingId,
        totalAmount: currentTotalAmount,
      });
      setPaymentStatus("processing");
      const paymentTimer = setTimeout(() => {
        setPaymentStatus("success");
      }, 3000);
      return () => clearTimeout(paymentTimer);
    } else {
      setError(
        "Không tìm thấy thông tin thanh toán hợp lệ. Vui lòng thử lại quy trình đặt vé."
      );
      setPaymentStatus("failed");
    }
  }, [location.state, bookingIdFromParams]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      navigate("/");
    }
  }, [isAuthenticated, authLoading, navigate]);

  useEffect(() => {
    if (paymentStatus === "success" && bookingInfo?.bookingId) {
      const successTimer = setTimeout(() => {
        navigate(`/ticket/${bookingInfo.bookingId}`);
      }, 1500);
      return () => clearTimeout(successTimer);
    }
  }, [paymentStatus, bookingInfo, navigate]);

  if (authLoading) {
    return (
      <Container className="text-center mt-5">
        <Spinner animation="border" />
      </Container>
    );
  }

  let content;
  if (paymentStatus === "pending" || !bookingInfo) {
    content = (
      <>
        <Spinner animation="grow" variant="info" className="me-2" />
        Đang tải thông tin thanh toán...
      </>
    );
  } else if (paymentStatus === "processing") {
    content = (
      <>
        <Spinner animation="border" variant="primary" className="me-2" />
        Đang xử lý thanh toán cho booking{" "}
        <strong>#{bookingInfo.bookingId}</strong>...
        <p className="mt-3">
          Tổng tiền:{" "}
          <strong className="text-danger">
            {bookingInfo.totalAmount.toLocaleString()} VNĐ
          </strong>
        </p>
        <Alert variant="info" className="mt-3">
          Đây là trang giả lập thanh toán.
        </Alert>
      </>
    );
  } else if (paymentStatus === "success") {
    content = (
      <Alert variant="success">
        <Alert.Heading>Thanh toán Thành công!</Alert.Heading>
        <p>
          Booking <strong>#{bookingInfo?.bookingId}</strong> của bạn đã được xác
          nhận.
        </p>
        <p>Đang chuẩn bị chuyển đến trang thông tin vé của bạn...</p>
      </Alert>
    );
  } else if (paymentStatus === "failed") {
    content = (
      <Alert variant="danger">
        <Alert.Heading>Thanh toán Thất bại!</Alert.Heading>
        <p>{error || "Đã có lỗi xảy ra trong quá trình thanh toán."}</p>
        <hr />
        <Button onClick={() => navigate("/")} variant="outline-danger">
          Quay về Trang chủ
        </Button>
      </Alert>
    );
  }

  return (
    <Container className="mt-5 d-flex justify-content-center">
      <Card
        className="shadow-lg p-4 text-center"
        style={{ minWidth: "400px", maxWidth: "600px" }}
      >
        <Card.Body>
          <Card.Title as="h2" className="mb-4">
            {paymentStatus === "success"
              ? "Giao dịch Hoàn tất"
              : "Thanh toán Đặt vé"}
          </Card.Title>
          {content}
        </Card.Body>
      </Card>
    </Container>
  );
};

export default PaymentPage;
