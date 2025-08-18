import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../components/Header';
import '../components/CleanerPanel.css';

const CleanerPanel = () => {
  const [user, setUser] = useState(null);
  const [toilets, setToilets] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  // Problem türleri
  const problemTypes = {
    1: "Tuvalet Kağıdı yok",
    2: "Sabun yok",
    3: "Peçete yok",
    4: "Çöp kutusu dolu",
    5: "Klozet kirli",
    6: "Diğer"
  };

  // Tuvaletlerin durumunu getir
  const fetchToiletsStatus = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/toilets/status');
      const data = await response.json();
      
      if (data.success) {
        setToilets(data.data);
      } else {
        console.error('Tuvalet durumu getirilemedi:', data.message);
      }
    } catch (error) {
      console.error('API hatası:', error);
    }
  };

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
      
      // Tuvalet durumlarını getir
      fetchToiletsStatus();
      
    } catch (error) {
      console.error('Kullanıcı verileri parse edilemedi:', error);
      navigate('/login');
    } finally {
      setLoading(false);
    }
  }, [navigate]);

  // Problemleri parse et
  const parseProblems = (problemsStr) => {
    try {
      if (!problemsStr) return [];
      const problemIds = JSON.parse(problemsStr);
      return problemIds.map(id => problemTypes[id]).filter(Boolean);
    } catch (error) {
      return [];
    }
  };

  // Durum rengi belirle
  const getStatusColor = (toilet) => {
    if (!toilet.last_rating) return '#6c757d'; // Gri - hiç kontrol edilmemiş
    if (toilet.has_problems) return '#dc3545'; // Kırmızı - problem var
    return '#28a745'; // Yeşil - temiz
  };

  // Durum metni belirle
  const getStatusText = (toilet) => {
    if (!toilet.last_rating) return 'Henüz kontrol edilmemiş';
    if (toilet.has_problems) return `${toilet.problem_count} problem tespit edildi`;
    return 'Temiz';
  };

  // Temizlik görevini başlat (ilk kez)
  const startCleaningTask = async (toiletId) => {
    try {
      const token = localStorage.getItem('authToken');
      const user = JSON.parse(localStorage.getItem('user'));
      
      const response = await fetch('http://localhost:8080/api/cleaning/start', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
          'X-User-ID': user.id.toString(),
          'X-User-Name': user.name
        },
        body: JSON.stringify({ toilet_id: toiletId })
      });
      
      const data = await response.json();
      
      if (data.success) {
        alert('Temizlik görevi alındı!');
        fetchToiletsStatus(); // Durumları yenile
      } else {
        alert(data.message);
      }
    } catch (error) {
      console.error('Temizlik görevi başlatılamadı:', error);
      alert('Temizlik görevi başlatılırken hata oluştu');
    }
  };

  // Temizlik görevini fiilen başlat (assigned -> in_progress)
  const beginCleaningTask = async (taskId) => {
    try {
      const token = localStorage.getItem('authToken');
      
      const response = await fetch(`http://localhost:8080/api/cleaning/begin/${taskId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });
      
      const data = await response.json();
      
      if (data.success) {
        alert('Temizlik başlatıldı!');
        fetchToiletsStatus(); // Durumları yenile
      } else {
        alert(data.message);
      }
    } catch (error) {
      console.error('Temizlik başlatılamadı:', error);
      alert('Temizlik başlatılırken hata oluştu');
    }
  };

  // Temizlik görevini tamamla
  const completeCleaningTask = async (taskId) => {
    try {
      const token = localStorage.getItem('authToken');
      
      const response = await fetch(`http://localhost:8080/api/cleaning/complete/${taskId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });
      
      const data = await response.json();
      
      if (data.success) {
        alert('Temizlik tamamlandı!');
        fetchToiletsStatus(); // Durumları yenile
      } else {
        alert(data.message);
      }
    } catch (error) {
      console.error('Temizlik tamamlanamadı:', error);
      alert('Temizlik tamamlanırken hata oluştu');
    }
  };

  // Buton durumunu belirle
  const getButtonConfig = (toilet) => {
    // Aktif temizlik görevi var mı?
    if (toilet.cleaning_task) {
      const task = toilet.cleaning_task;
      const currentUser = JSON.parse(localStorage.getItem('user'));
      
      // Bu kullanıcının görevi mi?
      if (task.cleaner_id === currentUser.id) {
        if (task.status === 'assigned') {
          return {
            text: 'Temizliğe Başla',
            action: () => beginCleaningTask(task.id),
            className: 'btn-warning'
          };
        } else if (task.status === 'in_progress') {
          return {
            text: 'Tamamlandı',
            action: () => completeCleaningTask(task.id),
            className: 'btn-success'
          };
        }
      } else {
        // Başka birinin görevi
        return {
          text: `${task.cleaner_name} Temizliyor`,
          action: null,
          className: 'btn-disabled'
        };
      }
    }
    
    // Tuvalet kirli mi?
    if (toilet.has_problems) {
      return {
        text: 'Temizliğe Başla',
        action: () => startCleaningTask(toilet.toilet.id),
        className: 'btn-primary'
      };
    }
    
    // Temizse değerlendirme butonu
    return {
      text: 'Değerlendir',
      action: () => navigate(`/rating?toilet=${toilet.toilet.id}`),
      className: 'btn-primary'
    };
  };

  // Zaman formatla
  const formatTime = (dateString) => {
    if (!dateString) return 'Henüz kontrol edilmemiş';
    const date = new Date(dateString);
    return date.toLocaleString('tr-TR');
  };

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

        <div className="toilets-section">
          <h2>Tuvalet Durumları</h2>
          <div className="toilets-grid">
            {toilets.map((toilet) => (
              <div key={toilet.toilet.id} className="toilet-card">
                <div className="toilet-header">
                  <h3>{toilet.toilet.name}</h3>
                  <div 
                    className="status-indicator"
                    style={{ backgroundColor: getStatusColor(toilet) }}
                  ></div>
                </div>
                
                <div className="toilet-info">
                  <p className="location">📍 {toilet.toilet.location}</p>
                  <p className="status-text">{getStatusText(toilet)}</p>
                  <p className="last-checked">
                    Son kontrol: {formatTime(toilet.last_checked)}
                  </p>
                  
                  {toilet.cleaning_task && (
                    <div className="cleaning-task-info">
                      <p className={`task-status task-${toilet.cleaning_task.status}`}>
                        🔧 {toilet.cleaning_task.status === 'assigned' ? 'Görev Alındı' : 
                           toilet.cleaning_task.status === 'in_progress' ? 'Temizlik Devam Ediyor' : 
                           'Temizlik Tamamlandı'}
                      </p>
                      <p className="cleaner-name">
                        👤 {toilet.cleaning_task.cleaner_name}
                      </p>
                      {toilet.cleaning_task.started_at && (
                        <p className="started-time">
                          🕐 Başlangıç: {formatTime(toilet.cleaning_task.started_at)}
                        </p>
                      )}
                    </div>
                  )}
                  
                  {toilet.last_rating && (
                    <div className="rating-info">
                      <p className="rating">
                        Puan: {toilet.last_rating.rating}/5 ⭐
                      </p>
                      
                      {toilet.has_problems && (
                        <div className="problems">
                          <strong>Tespit edilen problemler:</strong>
                          <ul>
                            {parseProblems(toilet.last_rating.problems).map((problem, index) => (
                              <li key={index}>{problem}</li>
                            ))}
                          </ul>
                          {toilet.last_rating.other_text && (
                            <p className="other-text">
                              <strong>Diğer:</strong> {toilet.last_rating.other_text}
                            </p>
                          )}
                        </div>
                      )}
                    </div>
                  )}
                </div>
                
                <div className="toilet-actions">
                  {(() => {
                    const buttonConfig = getButtonConfig(toilet);
                    return (
                      <button 
                        className={buttonConfig.className}
                        onClick={buttonConfig.action}
                        disabled={!buttonConfig.action}
                      >
                        {buttonConfig.text}
                      </button>
                    );
                  })()}
                </div>
              </div>
            ))}
          </div>
          
          <div className="refresh-section">
            <button 
              className="btn-secondary"
              onClick={fetchToiletsStatus}
            >
              Durumları Yenile
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default CleanerPanel;
