import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../components/Header';
import '../components/AdminPanel.css';

const AdminPanel = () => {
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
      
      // Admin kontrolü
      if (userData.role !== 'admin') {
        navigate('/cleaner'); // Admin değilse temizlikçi paneline yönlendir
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
      <div className="admin-panel">
        <div className="panel-header">
          <h1>Admin Paneli</h1>
          <p>Hoş geldiniz, {user.name}</p>
        </div>

        <div className="admin-stats">
          <div className="stat-card">
            <h3>Toplam Değerlendirme</h3>
            <div className="stat-number">--</div>
          </div>
          <div className="stat-card">
            <h3>Aktif Kullanıcı</h3>
            <div className="stat-number">--</div>
          </div>
          <div className="stat-card">
            <h3>Ortalama Puan</h3>
            <div className="stat-number">--</div>
          </div>
          <div className="stat-card">
            <h3>Sorunlu Alan</h3>
            <div className="stat-number">--</div>
          </div>
        </div>

        <div className="admin-sections">
          <div className="section-card">
            <h3>📊 Raporlar ve Analizler</h3>
            <p>Detaylı raporları görüntüleyin ve analiz edin</p>
            <div className="section-buttons">
              <button className="btn-primary">Günlük Rapor</button>
              <button className="btn-primary">Haftalık Rapor</button>
              <button className="btn-primary">Aylık Rapor</button>
            </div>
          </div>

          <div className="section-card">
            <h3>👥 Kullanıcı Yönetimi</h3>
            <p>Kullanıcıları ekleyin, düzenleyin ve yönetin</p>
            <div className="section-buttons">
              <button className="btn-secondary">Kullanıcı Listesi</button>
              <button className="btn-secondary">Yeni Kullanıcı Ekle</button>
              <button className="btn-secondary">Rol Düzenle</button>
            </div>
          </div>

          <div className="section-card">
            <h3>⚙️ Sistem Ayarları</h3>
            <p>Sistem konfigürasyonunu yönetin</p>
            <div className="section-buttons">
              <button className="btn-accent">Genel Ayarlar</button>
              <button className="btn-accent">Bildirim Ayarları</button>
              <button className="btn-accent">Backup/Restore</button>
            </div>
          </div>

          <div className="section-card">
            <h3>🏢 Alan Yönetimi</h3>
            <p>Temizlik alanlarını tanımlayın ve yönetin</p>
            <div className="section-buttons">
              <button className="btn-info">Alan Listesi</button>
              <button className="btn-info">Yeni Alan Ekle</button>
              <button className="btn-info">QR Kod Üret</button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default AdminPanel;
