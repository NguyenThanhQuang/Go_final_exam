import React, { useState, useEffect } from "react";
import {
  Container,
  Form,
  Button,
  Row,
  Col,
  Card,
  Spinner,
  Alert,
} from "react-bootstrap";
import type { Trip } from "../types/trip.types";
import { searchTrips, getAllTrips } from "../api/tripApi";
import { useAuth } from "../contexts/AuthContext";
import { useNavigate } from "react-router-dom";

interface HomePageProps {
  onShowAuthModal: () => void;
}

const HomePage: React.FC<HomePageProps> = ({ onShowAuthModal }) => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [date, setDate] = useState("");

  const [displayedTrips, setDisplayedTrips] = useState<Trip[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchInitialTrips = async () => {
      setLoading(true);
      setError(null);
      try {
        const results = await getAllTrips();
        if (results && results.length > 0) {
          setDisplayedTrips(results);
        } else {
          setDisplayedTrips([]);
        }
      } catch (err) {
        setError("Không thể tải danh sách chuyến đi ban đầu.");
        setDisplayedTrips([]);
      } finally {
        setLoading(false);
      }
    };
    fetchInitialTrips();
  }, []);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!from || !to || !date) {
      setError("Vui lòng nhập đầy đủ thông tin điểm đi, điểm đến và ngày đi.");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const results = await searchTrips({ from, to, date });
      setDisplayedTrips(results);
      if (results.length === 0) {
        setError("Không tìm thấy chuyến đi nào phù hợp với tìm kiếm của bạn.");
      }
    } catch (err) {
      setError("Lỗi khi tìm kiếm. Vui lòng thử lại.");
    } finally {
      setLoading(false);
    }
  };

  const handleBookTripClick = (tripId: string) => {
    if (isAuthenticated) {
      navigate(`/trips/${tripId}`);
    } else {
      onShowAuthModal();
    }
  };

  return (
    <Container fluid className="mt-4 px-md-5">
      <Container>
        <Card className="p-4 mb-4 shadow">
          <h2 className="mb-3 text-center">Tìm Vé Xe Rẻ - Đặt Vé Dễ Dàng</h2>
          <Form onSubmit={handleSearch}>
            <Row className="align-items-end">
              <Col md={3}>
                <Form.Group>
                  <Form.Label>Điểm đi</Form.Label>
                  <Form.Control
                    type="text"
                    value={from}
                    onChange={(e) => setFrom(e.target.value)}
                    placeholder="VD: Hà Nội"
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group>
                  <Form.Label>Điểm đến</Form.Label>
                  <Form.Control
                    type="text"
                    value={to}
                    onChange={(e) => setTo(e.target.value)}
                    placeholder="VD: Sài Gòn"
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Form.Group>
                  <Form.Label>Ngày đi</Form.Label>
                  <Form.Control
                    type="date"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                  />
                </Form.Group>
              </Col>
              <Col md={3}>
                <Button type="submit" className="w-100" disabled={loading}>
                  {loading ? (
                    <Spinner as="span" animation="border" size="sm" />
                  ) : (
                    "Tìm chuyến"
                  )}
                </Button>
              </Col>
            </Row>
          </Form>
        </Card>
      </Container>

      {error && (
        <Alert
          variant="danger"
          onClose={() => setError(null)}
          dismissible
          className="mt-3"
        >
          {error}
        </Alert>
      )}

      {loading && (
        <div className="text-center my-5">
          <Spinner
            animation="border"
            variant="primary"
            style={{ width: "3rem", height: "3rem" }}
          />
          <p className="mt-2">Đang tải dữ liệu...</p>
        </div>
      )}

      {!loading && displayedTrips.length === 0 && !error && (
        <Alert variant="info" className="mt-3">
          Hiện chưa có chuyến đi nào được hiển thị. Hãy thử tìm kiếm.
        </Alert>
      )}

      {!loading && displayedTrips.length > 0 && (
        <Row xs={1} md={2} lg={3} xl={4} className="g-4 mt-3">
          {displayedTrips.map((trip) => (
            <Col key={trip.id}>
              <Card className="h-100 shadow-sm">
                <Card.Body className="d-flex flex-column">
                  <Card.Title className="text-primary">
                    {trip.companyName || "Nhà xe ABC"}
                  </Card.Title>
                  <Card.Text>
                    <strong>{trip.route.from.name}</strong>{" "}
                    <i className="bi bi-arrow-right"></i>{" "}
                    <strong>{trip.route.to.name}</strong>
                  </Card.Text>
                  <Card.Text>
                    Khởi hành:{" "}
                    {new Date(trip.departureTime).toLocaleTimeString([], {
                      hour: "2-digit",
                      minute: "2-digit",
                    })}{" "}
                    - {new Date(trip.departureTime).toLocaleDateString()}
                  </Card.Text>
                  <div className="mt-auto d-flex justify-content-between align-items-center">
                    <div>
                      <h5 className="mb-0 text-danger">
                        {trip.price.toLocaleString()} VNĐ
                      </h5>
                      <small>Còn {trip.availableSeats} ghế</small>
                    </div>
                    <Button
                      variant="success"
                      onClick={() => handleBookTripClick(trip.id)}
                    >
                      Đặt vé
                    </Button>
                  </div>
                </Card.Body>
              </Card>
            </Col>
          ))}
        </Row>
      )}
    </Container>
  );
};

export default HomePage;
