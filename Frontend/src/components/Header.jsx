import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Header.css';

const Header = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    // localStorage'dan kullanıcı bilgilerini kontrol et
    const savedUser = localStorage.getItem('user');
    const token = localStorage.getItem('authToken');
    
    if (savedUser && token) {
      try {
        const userData = JSON.parse(savedUser);
        setUser(userData);
        setIsLoggedIn(true);
      } catch (error) {
        console.error('Kullanıcı verileri parse edilemedi:', error);
        // Hatalı veri varsa temizle
        localStorage.removeItem('user');
        localStorage.removeItem('authToken');
      }
    }
  }, []);

  const handleLogoClick = () => {
    navigate('/');
  };

  const handleLogin = () => {
    navigate('/login');
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
    setIsLoggedIn(false);
    setUser(null);
    // Ana sayfaya yönlendir
    navigate('/');
  };

  return (
    <header className="header">
      <div className="header-container">
        <div className="logo" onClick={handleLogoClick} style={{ cursor: 'pointer' }}>
          <h1>Temizlik Takip Sistemi</h1>
        </div>
        <div className="header-actions">
          {isLoggedIn ? (
            <div className="user-section">
              <span className="welcome-text">
                Hoş geldiniz, {user?.name}
              </span>
              <button className="logout-btn" onClick={handleLogout}>
                Çıkış Yap
              </button>
            </div>
          ) : (
            <button className="login-btn" onClick={handleLogin}>
              Giriş Yap
            </button>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;