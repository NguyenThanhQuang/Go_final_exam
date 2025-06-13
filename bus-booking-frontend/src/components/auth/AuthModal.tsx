import React, { useState } from "react";
import { Modal, Button, Form, Alert, Spinner } from "react-bootstrap";
import { useAuth } from "../../contexts/AuthContext";
import {
  loginUser,
  registerUser,
  type LoginPayload,
  type RegisterPayload,
} from "../../api/authApi";

const LoginForm: React.FC<{
  onSwitchToRegister: () => void;
  handleCloseModal: () => void;
}> = ({ onSwitchToRegister, handleCloseModal }) => {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const payload: LoginPayload = { email, password };
      const response = await loginUser(payload);
      if (response.token) {
        login(response.token);
        handleCloseModal();
      } else {
        setError(
          response["thông báo"] || "Đăng nhập thất bại, không nhận được token."
        );
      }
    } catch (err: any) {
      setError(
        err.response?.data?.lỗi ||
          err.message ||
          "Lỗi đăng nhập không xác định."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form onSubmit={handleSubmit}>
      {error && <Alert variant="danger">{error}</Alert>}
      <Form.Group className="mb-3">
        <Form.Label>Email</Form.Label>
        <Form.Control
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Mật khẩu</Form.Label>
        <Form.Control
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </Form.Group>
      <Button
        variant="primary"
        type="submit"
        className="w-100"
        disabled={loading}
      >
        {loading ? <Spinner animation="border" size="sm" /> : "Đăng nhập"}
      </Button>
      <div className="mt-3 text-center">
        Chưa có tài khoản?{" "}
        <Button variant="link" onClick={onSwitchToRegister} disabled={loading}>
          Đăng ký ngay
        </Button>
      </div>
    </Form>
  );
};

const RegisterForm: React.FC<{
  onSwitchToLogin: () => void;
  handleSuccessfulRegister: () => void;
}> = ({ onSwitchToLogin, handleSuccessfulRegister }) => {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccessMessage(null);
    setLoading(true);
    try {
      const payload: RegisterPayload = { name, email, phone, password };
      const response = await registerUser(payload);
      setSuccessMessage(
        response["thông báo"] || "Đăng ký thành công! Vui lòng đăng nhập."
      );
      setTimeout(() => {
        handleSuccessfulRegister();
      }, 1500);
    } catch (err: any) {
      setError(
        err.response?.data?.lỗi || err.message || "Lỗi đăng ký không xác định."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form onSubmit={handleSubmit}>
      {error && <Alert variant="danger">{error}</Alert>}
      {successMessage && <Alert variant="success">{successMessage}</Alert>}
      <Form.Group className="mb-3">
        <Form.Label>Họ và tên</Form.Label>
        <Form.Control
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          disabled={loading}
        />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Email</Form.Label>
        <Form.Control
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          disabled={loading}
        />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Số điện thoại</Form.Label>
        <Form.Control
          type="tel"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          required
          disabled={loading}
        />
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Mật khẩu</Form.Label>
        <Form.Control
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={6}
          disabled={loading}
        />
      </Form.Group>
      <Button
        variant="primary"
        type="submit"
        className="w-100"
        disabled={loading || !!successMessage}
      >
        {loading ? <Spinner animation="border" size="sm" /> : "Đăng ký"}
      </Button>
      <div className="mt-3 text-center">
        Đã có tài khoản?{" "}
        <Button variant="link" onClick={onSwitchToLogin} disabled={loading}>
          Đăng nhập
        </Button>
      </div>
    </Form>
  );
};

interface AuthModalProps {
  show: boolean;
  handleClose: () => void;
  mode: "login" | "register";
  setMode: (mode: "login" | "register") => void;
}

const AuthModal: React.FC<AuthModalProps> = ({
  show,
  handleClose,
  mode,
  setMode,
}) => {
  const handleRegisterSuccess = () => {
    setMode("login");
  };

  return (
    <Modal show={show} onHide={handleClose} centered>
      <Modal.Header closeButton>
        <Modal.Title>
          {mode === "login" ? "Đăng nhập" : "Đăng ký Tài khoản"}
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {mode === "login" ? (
          <LoginForm
            onSwitchToRegister={() => setMode("register")}
            handleCloseModal={handleClose}
          />
        ) : (
          <RegisterForm
            onSwitchToLogin={() => setMode("login")}
            handleSuccessfulRegister={handleRegisterSuccess}
          />
        )}
      </Modal.Body>
    </Modal>
  );
};

export default AuthModal;
