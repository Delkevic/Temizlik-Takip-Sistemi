import React from 'react';
import '../components/HomePage.css';
import Header from '../components/Header';

const HomePage = () => {
  return (
    <>
      <Header />
      <div className="home-page">
        <div className="guest-content">
          <h1>Temizlik Takip Sistemi</h1>
          <h2>Kaliteli temizlik hizmeti için değerlendirmeleriniz önemli.</h2>
          <p>Bu sistem, temizlik kalitesini takip etmek ve sürekli iyileştirme sağlamak için tasarlanmıştır.</p>
          <p>Sistem kullanımı için lütfen giriş yapınız.</p>
        </div>
      </div>
    </>
  );
}

export default HomePage