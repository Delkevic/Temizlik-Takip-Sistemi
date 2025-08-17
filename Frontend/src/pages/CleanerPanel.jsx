import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../components/Header';
import './CleanerPanel.css';

const CleanerPanel = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    // Kullanıcı kontrolü
    const savedUser = localStorage.getItem('user');
    const token = localStorage.getItem('authToken');
    
    if (!savedUser || !token) {
      navigate('/login');
      return;
    }
    
    try {
      const userData = JSON.parse(savedUser);
      
      // Temizlikçi kontrolü
      if (userData.role !== 'temizlikci') {
        navigate('/admin'); // Temizlikçi değilse admin paneline yönlendir
        return;
      }
      
      setUser(userData);
    } catch (error) {
      console.error('Kullanıcı verileri parse edilemedi:', error);
      navigate('/login');
    } finally {
      setLoading(false);
    }
  }, [navigate]);

  if (loading) {
    return <div className="loading">Yükleniyor...</div>;
  }

  if (!user) {
    return null;
  }

  return (
    <>
      <Header />
      <div className="cleaner-panel">
        <div className="panel-header">
          <h1>Temizlik Görevlisi Paneli</h1>
          <p>Hoş geldiniz, {user.name}</p>
        </div>

        <div className="cleaner-stats">
          <div className="stat-card">
            <h3>Bugünkü Görevler</h3>
            <div className="stat-number">--</div>
          </div>
          <div className="stat-card">
            <h3>Tamamlanan</h3>
            <div className="stat-number">--</div>
          </div>
          <div className="stat-card">
            <h3>Ortalama Puanım</h3>
            <div className="stat-number">--</div>
          </div>
          <div className="stat-card">
            <h3>Bu Ay</h3>
            <div className="stat-number">-- Görev</div>
          </div>
        </div>

        <div className="cleaner-sections">
          <div className="section-card">
            <h3>📋 Günlük Görevlerim</h3>
            <p>Bugün için atanmış görevlerinizi görüntüleyin</p>
            <div className="section-buttons">
              <button className="btn-primary">Aktif Görevler</button>
              <button className="btn-primary">Tamamlanan Görevler</button>
              <button className="btn-primary">Bekleyen Görevler</button>
            </div>
          </div>

          <div className="section-card">
            <h3>⭐ Puanlama ve Değerlendirme</h3>
            <p>Temizlik kalitesi değerlendirmesi yapın</p>
            <div className="section-buttons">
              <button className="btn-secondary">QR Kod Tara</button>
              <button className="btn-secondary">Manuel Değerlendirme</button>
              <button className="btn-secondary">Sorun Bildir</button>
            </div>
          </div>

          <div className="section-card">
            <h3>📊 Performans Takibi</h3>
            <p>Kendi performansınızı takip edin</p>
            <div className="section-buttons">
              <button className="btn-accent">Haftalık Özet</button>
              <button className="btn-accent">Aylık Rapor</button>
              <button className="btn-accent">Hedef Takibi</button>
            </div>
          </div>

          <div className="section-card">
            <h3>🔔 Bildirimler ve Mesajlar</h3>
            <p>Sistem bildirimleri ve mesajlarınız</p>
            <div className="section-buttons">
              <button className="btn-info">Yeni Bildirimler</button>
              <button className="btn-info">Geçmiş Mesajlar</button>
              <button className="btn-info">Duyurular</button>
            </div>
          </div>

          <div className="section-card">
            <h3>🛠️ Araçlar ve Kaynaklar</h3>
            <p>İş süreçlerinizde size yardımcı olacak araçlar</p>
            <div className="section-buttons">
              <button className="btn-warning">Temizlik Rehberi</button>
              <button className="btn-warning">Malzeme Talep</button>
              <button className="btn-warning">Yardım & Destek</button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default CleanerPanel;
