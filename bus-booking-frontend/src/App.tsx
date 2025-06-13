import React, { useState, useEffect, type JSX } from 'react';
import { BrowserRouter as Router, Routes, Route, Link, Navigate } from 'react-router-dom';
import { Navbar, Container, Nav, Button, NavDropdown, Spinner } from 'react-bootstrap';
import HomePage from './pages/HomePage';
import TripDetailsPage from './pages/TripDetailsPage';
import BookingConfirmationPage from './pages/BookingConfirmationPage';
import PaymentPage from './pages/PaymentPage';
import TicketDetailsPage from './pages/TicketDetailsPage';
import BookingHistoryPage from './pages/BookingHistoryPage';
import AuthModal from './components/auth/AuthModal';
import { useAuth } from './contexts/AuthContext';

const ProtectedRoute: React.FC<{ children: JSX.Element }> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <Container className="text-center mt-5"><Spinner animation="border" /></Container>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/?showLogin=true" replace />;
  }

  return children;
};


function App() {
  const { isAuthenticated, logout, isLoading } = useAuth();
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [authMode, setAuthMode] = useState<'login' | 'register'>('login');

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    if (params.get('showLogin') === 'true' && !isAuthenticated && !isLoading) {
      handleShowLogin();
    }
  }, [isAuthenticated, isLoading]);


  const handleShowLogin = () => {
    setAuthMode('login');
    setShowAuthModal(true);
  };

  const handleShowRegister = () => {
    setAuthMode('register');
    setShowAuthModal(true);
  };

  const handleCloseAuthModal = () => setShowAuthModal(false);

  const handleLogout = () => {
    logout();
  };

  if (isLoading) {
    return (
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
            <Spinner animation="border" variant="primary" />
        </div>
    );
  }

  return (
    <Router>
      <Navbar bg="dark" variant="dark" expand="lg" sticky="top">
        <Container>
          <Navbar.Brand as={Link} to="/">Vé Xe Khách</Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="ms-auto align-items-center">
              {isAuthenticated ? (
                <NavDropdown title="Tài khoản" id="user-dropdown">
                  <NavDropdown.Item as={Link} to="/my-bookings">Vé của tôi</NavDropdown.Item>
                  <NavDropdown.Divider />
                  <NavDropdown.Item onClick={handleLogout}>Đăng xuất</NavDropdown.Item>
                </NavDropdown>
              ) : (
                <>
                  <Button variant="outline-light" onClick={handleShowLogin} className="me-2">
                    Đăng nhập
                  </Button>
                  <Button variant="light" onClick={handleShowRegister}>
                    Đăng ký
                  </Button>
                </>
              )}
            </Nav>
          </Navbar.Collapse>
        </Container>
      </Navbar>

      <Container fluid className="mt-3">
        <Routes>
          <Route path="/" element={<HomePage onShowAuthModal={handleShowLogin} />} />
          <Route path="/trips/:tripId" element={<TripDetailsPage />} />
          <Route path="/booking/confirm" element={
            <ProtectedRoute><BookingConfirmationPage /></ProtectedRoute>
          } />
          <Route path="/payment/:bookingId" element={
            <ProtectedRoute><PaymentPage /></ProtectedRoute>
          } />
          <Route path="/ticket/:bookingId" element={
            <ProtectedRoute><TicketDetailsPage /></ProtectedRoute>
          } />
          <Route path="/my-bookings" element={
            <ProtectedRoute>
              <BookingHistoryPage />
            </ProtectedRoute>
          } />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </Container>

      {!isAuthenticated && (
          <AuthModal 
            show={showAuthModal} 
            handleClose={handleCloseAuthModal} 
            mode={authMode} 
            setMode={setAuthMode} 
          />
      )}
    </Router>
  );
}

export default App;