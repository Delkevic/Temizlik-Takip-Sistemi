import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../components/HomePage.css';
import Header from '../components/Header';

const HomePage = () => {
  const navigate = useNavigate();

  return (
    <>
      <Header />
      <div className="home-page">
        <div className="guest-content">
          <h1>Temizlik Takip Sistemi</h1>
          <h2>Kaliteli temizlik hizmeti için değerlendirmeleriniz önemli.</h2>
          <p>Bu sistem, temizlik kalitesini takip etmek ve sürekli iyileştirme sağlamak için tasarlanmıştır.</p>
          
          <div className="action-buttons">
            <button 
              className="rating-button"
              onClick={() => navigate('/rating')}
            >
              Tuvalet Değerlendirme Yap
            </button>
            <p className="login-text">Personel girişi için lütfen giriş yapınız.</p>
          </div>
        </div>
      </div>
    </>
  );
}

export default HomePage