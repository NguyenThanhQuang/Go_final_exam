import React, { useState } from "react";
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
import { searchTrips } from "../api/tripApi";

const HomePage: React.FC = () => {
  const [from, setFrom] = useState("TP. Hồ Chí Minh");
  const [to, setTo] = useState("Đà Lạt");
  const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
  const [trips, setTrips] = useState<Trip[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setTrips([]);

    try {
      const results = await searchTrips({ from, to, date });
      setTrips(results);
    } catch (err) {
      setError(
        "Không thể kết nối đến máy chủ hoặc đã có lỗi xảy ra. Vui lòng thử lại."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container className="mt-4">
      <Card className="p-4">
        <h2 className="mb-4">Tìm kiếm chuyến đi</h2>
        <Form onSubmit={handleSearch}>
          <Row>
            <Col md={4}>
              <Form.Group className="mb-3">
                <Form.Label>Điểm đi</Form.Label>
                <Form.Control
                  type="text"
                  value={from}
                  onChange={(e) => setFrom(e.target.value)}
                  required
                />
              </Form.Group>
            </Col>
            <Col md={4}>
              <Form.Group className="mb-3">
                <Form.Label>Điểm đến</Form.Label>
                <Form.Control
                  type="text"
                  value={to}
                  onChange={(e) => setTo(e.target.value)}
                  required
                />
              </Form.Group>
            </Col>
            <Col md={3}>
              <Form.Group className="mb-3">
                <Form.Label>Ngày đi</Form.Label>
                <Form.Control
                  type="date"
                  value={date}
                  onChange={(e) => setDate(e.target.value)}
                  required
                />
              </Form.Group>
            </Col>
            <Col md={1} className="d-flex align-items-end mb-3">
              <Button type="submit" className="w-100" disabled={loading}>
                {loading ? (
                  <Spinner as="span" animation="border" size="sm" />
                ) : (
                  "Tìm"
                )}
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>

      <div className="mt-4">
        {error && <Alert variant="danger">{error}</Alert>}

        {loading ? (
          <div className="text-center">
            <Spinner animation="border" role="status">
              <span className="visually-hidden">Đang tải...</span>
            </Spinner>
          </div>
        ) : trips.length > 0 ? (
          trips.map((trip) => (
            <Card key={trip.id} className="mb-3">
              <Card.Body>
                <Row>
                  <Col md={8}>
                    <Card.Title>{trip.companyName}</Card.Title>
                    <Card.Text>
                      <strong>{trip.route.from.name}</strong> →{" "}
                      <strong>{trip.route.to.name}</strong>
                    </Card.Text>
                    <Card.Text>
                      Khởi hành:{" "}
                      {new Date(trip.departureTime).toLocaleTimeString()} -{" "}
                      {new Date(trip.departureTime).toLocaleDateString()}
                    </Card.Text>
                  </Col>
                  <Col md={4} className="text-md-end">
                    <h4 className="text-danger">
                      {trip.price.toLocaleString()} VNĐ
                    </h4>
                    <div>
                      Còn <strong>{trip.availableSeats}</strong> ghế trống
                    </div>
                    <Button variant="primary" className="mt-2">
                      Chọn chuyến
                    </Button>
                  </Col>
                </Row>
              </Card.Body>
            </Card>
          ))
        ) : (
          !error && (
            <Alert variant="info">Không có chuyến đi nào được tìm thấy.</Alert>
          )
        )}
      </div>
    </Container>
  );
};

export default HomePage;
