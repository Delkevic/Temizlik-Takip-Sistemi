import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Header from '../components/Header';
import '../components/AdminPanel.css';

const AdminPanel = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [toiletStatuses, setToiletStatuses] = useState([]);
  const [loadingToilets, setLoadingToilets] = useState(false);
  const [showAddUser, setShowAddUser] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [activeTab, setActiveTab] = useState('toilets'); // toilets, users, stats
  const [stats, setStats] = useState(null);
  const [loadingStats, setLoadingStats] = useState(false);
  const [newUser, setNewUser] = useState({
    username: '',
    password: '',
    name: '',
    role: 'temizlikci'
  });
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
      fetchUsers();
      fetchToiletStatuses();
      fetchStats();
    } catch (error) {
      console.error('Kullanıcı verileri parse edilemedi:', error);
      navigate('/login');
    } finally {
      setLoading(false);
    }
  }, [navigate]);

  const fetchUsers = async () => {
    setLoadingUsers(true);
    try {
      const response = await fetch('http://localhost:8080/api/admin/users');
      const data = await response.json();
      if (data.success) {
        setUsers(data.users || []);
      } else {
        console.error('Kullanıcılar getirilemedi:', data.message);
      }
    } catch (error) {
      console.error('Kullanıcılar getirilirken hata oluştu:', error);
    } finally {
      setLoadingUsers(false);
    }
  };

  const fetchToiletStatuses = async () => {
    setLoadingToilets(true);
    try {
      const response = await fetch('http://localhost:8080/api/toilets/status');
      const data = await response.json();
      if (data.success) {
        setToiletStatuses(data.data || []);
      } else {
        console.error('Tuvalet durumları getirilemedi:', data.message);
      }
    } catch (error) {
      console.error('Tuvalet durumları getirilirken hata oluştu:', error);
    } finally {
      setLoadingToilets(false);
    }
  };

  const fetchStats = async () => {
    setLoadingStats(true);
    try {
      const response = await fetch('http://localhost:8080/api/admin/stats');
      const data = await response.json();
      if (data.success) {
        setStats(data);
      } else {
        console.error('İstatistikler getirilemedi:', data.message);
      }
    } catch (error) {
      console.error('İstatistikler getirilirken hata oluştu:', error);
    } finally {
      setLoadingStats(false);
    }
  };

  const handleAddUser = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:8080/api/admin/users', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newUser),
      });

      const data = await response.json();
      if (data.success) {
        alert('Kullanıcı başarıyla eklendi!');
        setNewUser({ username: '', password: '', name: '', role: 'temizlikci' });
        setShowAddUser(false);
        fetchUsers();
      } else {
        alert('Hata: ' + data.message);
      }
    } catch (error) {
      console.error('Kullanıcı eklenirken hata oluştu:', error);
      alert('Kullanıcı eklenirken hata oluştu!');
    }
  };

  const handleEditUser = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`http://localhost:8080/api/admin/users/${editingUser.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: editingUser.username,
          name: editingUser.name,
          role: editingUser.role,
          ...(editingUser.password && { password: editingUser.password })
        }),
      });

      const data = await response.json();
      if (data.success) {
        alert('Kullanıcı başarıyla güncellendi!');
        setEditingUser(null);
        fetchUsers();
      } else {
        alert('Hata: ' + data.message);
      }
    } catch (error) {
      console.error('Kullanıcı güncellenirken hata oluştu:', error);
      alert('Kullanıcı güncellenirken hata oluştu!');
    }
  };

  const handleDeleteUser = async (userId) => {
    if (!confirm('Bu kullanıcıyı silmek istediğinizden emin misiniz?')) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/admin/users/${userId}`, {
        method: 'DELETE',
      });

      const data = await response.json();
      if (data.success) {
        alert('Kullanıcı başarıyla silindi!');
        fetchUsers();
      } else {
        alert('Hata: ' + data.message);
      }
    } catch (error) {
      console.error('Kullanıcı silinirken hata oluştu:', error);
      alert('Kullanıcı silinirken hata oluştu!');
    }
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
      <div className="admin-panel">
        <div className="panel-header">
          <h1>Admin Paneli</h1>
          <p>Hoş geldiniz, {user.name}</p>
        </div>

        {/* Tab Navigation */}
        <div className="tab-navigation">
          <button 
            className={`tab-btn ${activeTab === 'toilets' ? 'active' : ''}`}
            onClick={() => setActiveTab('toilets')}
          >
            Tuvalet Durumları
          </button>
          <button 
            className={`tab-btn ${activeTab === 'users' ? 'active' : ''}`}
            onClick={() => setActiveTab('users')}
          >
            Temizlikçi Yönetimi
          </button>
          <button 
            className={`tab-btn ${activeTab === 'stats' ? 'active' : ''}`}
            onClick={() => setActiveTab('stats')}
          >
            İstatistikler
          </button>
        </div>

        {/* Tuvalet Durumları Sekmesi */}
        {activeTab === 'toilets' && (
          <div className="toilet-status-section">
            <div className="section-header">
              <h2>Tuvalet Durumları</h2>
              <button 
                className="btn-secondary" 
                onClick={fetchToiletStatuses}
              >
                Yenile
              </button>
            </div>

            {loadingToilets ? (
              <div className="loading">Tuvalet durumları yükleniyor...</div>
            ) : (
              <div className="toilet-status-grid">
                {toiletStatuses.map(status => (
                  <div key={status.toilet.id} className="toilet-status-card">
                    <div className="toilet-info">
                      <h3>{status.toilet.name}</h3>
                      <p className="location">{status.toilet.location}</p>
                    </div>
                    
                    <div className="status-details">
                      <div className="rating-info">
                        <div className="average-rating">
                          <span className="label">Ortalama Puan:</span>
                          <span className={`rating-value ${status.average_rating >= 4 ? 'good' : status.average_rating >= 3 ? 'medium' : 'poor'}`}>
                            {status.total_ratings > 0 ? status.average_rating.toFixed(1) : 'Puan yok'}
                          </span>
                          {status.total_ratings > 0 && (
                            <span className="rating-count">({status.total_ratings} değerlendirme)</span>
                          )}
                        </div>
                      </div>

                      <div className="problem-status">
                        <span className={`status-badge ${status.has_problems ? 'problems' : 'no-problems'}`}>
                          {status.has_problems ? `${status.problem_count} Problem` : 'Problem Yok'}
                        </span>
                      </div>

                      <div className="cleaning-status">
                        {status.cleaning_task ? (
                          <span className={`status-badge ${status.cleaning_task.status}`}>
                            {status.cleaning_task.status === 'assigned' && 'Temizlik Atandı'}
                            {status.cleaning_task.status === 'in_progress' && 'Temizlik Devam Ediyor'}
                          </span>
                        ) : (
                          <span className="status-badge idle">Temizlik Gerekmiyor</span>
                        )}
                      </div>

                      {status.last_checked && (
                        <div className="last-checked">
                          <span className="label">Son Kontrol:</span>
                          <span className="time">
                            {new Date(status.last_checked).toLocaleString('tr-TR')}
                          </span>
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Temizlikçi Yönetimi Sekmesi */}
        {activeTab === 'users' && (
          <div className="user-management">
            <div className="section-header">
              <h2>Temizlikçi Yönetimi</h2>
              <button 
                className="btn-primary" 
                onClick={() => setShowAddUser(true)}
              >
                Yeni Temizlikçi Ekle
              </button>
            </div>

            {loadingUsers ? (
              <div className="loading">Kullanıcılar yükleniyor...</div>
            ) : (
              <div className="users-table">
                <table>
                  <thead>
                    <tr>
                      <th>ID</th>
                      <th>Kullanıcı Adı</th>
                      <th>İsim</th>
                      <th>Rol</th>
                      <th>Durum</th>
                      <th>İşlemler</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map(u => (
                      <tr key={u.id}>
                        <td>{u.id}</td>
                        <td>{u.username}</td>
                        <td>{u.name}</td>
                        <td>{u.role}</td>
                        <td>
                          <span className={`status ${u.is_active ? 'active' : 'inactive'}`}>
                            {u.is_active ? 'Aktif' : 'Pasif'}
                          </span>
                        </td>
                        <td>
                          <button 
                            className="btn-secondary" 
                            onClick={() => setEditingUser({...u, password: ''})}
                          >
                            Düzenle
                          </button>
                          <button 
                            className="btn-danger" 
                            onClick={() => handleDeleteUser(u.id)}
                          >
                            Sil
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        )}

        {/* İstatistikler Sekmesi */}
        {activeTab === 'stats' && (
          <div className="stats-section">
            <div className="section-header">
              <h2>İstatistikler</h2>
              <button 
                className="btn-secondary" 
                onClick={fetchStats}
              >
                Yenile
              </button>
            </div>

            {loadingStats ? (
              <div className="loading">İstatistikler yükleniyor...</div>
            ) : stats ? (
              <div className="stats-container">
                {/* Sistem İstatistikleri */}
                <div className="system-stats">
                  <h3>Sistem Özeti</h3>
                  <div className="stats-grid">
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.total_toilets}</div>
                      <div className="stat-label">Toplam Tuvalet</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.active_toilets}</div>
                      <div className="stat-label">Aktif Tuvalet</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.toilets_with_problems}</div>
                      <div className="stat-label">Problem Olan Tuvalet</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.total_cleaners}</div>
                      <div className="stat-label">Toplam Temizlikçi</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.active_cleaners}</div>
                      <div className="stat-label">Aktif Temizlikçi</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.total_ratings}</div>
                      <div className="stat-label">Toplam Değerlendirme</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">
                        {stats.system_stats.average_rating ? stats.system_stats.average_rating.toFixed(1) : '0.0'}
                      </div>
                      <div className="stat-label">Ortalama Puan</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.completed_tasks_today}</div>
                      <div className="stat-label">Bugün Tamamlanan</div>
                    </div>
                    <div className="stat-card">
                      <div className="stat-value">{stats.system_stats.ongoing_tasks}</div>
                      <div className="stat-label">Devam Eden Görev</div>
                    </div>
                  </div>
                </div>

                {/* Temizlikçi İstatistikleri */}
                <div className="cleaner-stats">
                  <h3>Temizlikçi Performansları</h3>
                  <div className="cleaner-stats-table">
                    <table>
                      <thead>
                        <tr>
                          <th>İsim</th>
                          <th>Durum</th>
                          <th>Toplam Temizlik</th>
                          <th>Ortalama Süre</th>
                          <th>En Hızlı</th>
                          <th>En Yavaş</th>
                          <th>Bu Hafta</th>
                          <th>Bu Ay</th>
                          <th>Devam Eden</th>
                        </tr>
                      </thead>
                      <tbody>
                        {stats.cleaner_stats.map(cleaner => (
                          <tr key={cleaner.cleaner_id}>
                            <td>{cleaner.cleaner_name}</td>
                            <td>
                              <span className={`status ${cleaner.is_active ? 'active' : 'inactive'}`}>
                                {cleaner.is_active ? 'Aktif' : 'Pasif'}
                              </span>
                            </td>
                            <td>{cleaner.total_completed_tasks}</td>
                            <td>
                              {cleaner.average_cleaning_time >= 0 
                                ? `${Math.round(cleaner.average_cleaning_time)} dk` 
                                : '-'}
                            </td>
                            <td>
                              {cleaner.fastest_cleaning_time >= 0 
                                ? `${Math.round(cleaner.fastest_cleaning_time)} dk` 
                                : '-'}
                            </td>
                            <td>
                              {cleaner.slowest_cleaning_time >= 0 
                                ? `${Math.round(cleaner.slowest_cleaning_time)} dk` 
                                : '-'}
                            </td>
                            <td>{cleaner.last_week_tasks}</td>
                            <td>{cleaner.last_month_tasks}</td>
                            <td>{cleaner.ongoing_tasks}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            ) : (
              <div>İstatistikler yüklenemedi</div>
            )}
          </div>
        )}

        {/* Modal'lar */}

        {/* Yeni Kullanıcı Ekleme Modal */}
        {showAddUser && (
          <div className="modal-overlay">
            <div className="modal">
              <h3>Yeni Temizlikçi Ekle</h3>
              <form onSubmit={handleAddUser}>
                <div className="form-group">
                  <label>Kullanıcı Adı:</label>
                  <input
                    type="text"
                    value={newUser.username}
                    onChange={(e) => setNewUser({...newUser, username: e.target.value})}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>İsim:</label>
                  <input
                    type="text"
                    value={newUser.name}
                    onChange={(e) => setNewUser({...newUser, name: e.target.value})}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Şifre:</label>
                  <input
                    type="password"
                    value={newUser.password}
                    onChange={(e) => setNewUser({...newUser, password: e.target.value})}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Rol:</label>
                  <select
                    value={newUser.role}
                    onChange={(e) => setNewUser({...newUser, role: e.target.value})}
                  >
                    <option value="temizlikci">Temizlikçi</option>
                    <option value="admin">Admin</option>
                  </select>
                </div>
                <div className="modal-buttons">
                  <button type="submit" className="btn-primary">Ekle</button>
                  <button 
                    type="button" 
                    className="btn-secondary" 
                    onClick={() => setShowAddUser(false)}
                  >
                    İptal
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* Kullanıcı Düzenleme Modal */}
        {editingUser && (
          <div className="modal-overlay">
            <div className="modal">
              <h3>Kullanıcı Düzenle</h3>
              <form onSubmit={handleEditUser}>
                <div className="form-group">
                  <label>Kullanıcı Adı:</label>
                  <input
                    type="text"
                    value={editingUser.username}
                    onChange={(e) => setEditingUser({...editingUser, username: e.target.value})}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>İsim:</label>
                  <input
                    type="text"
                    value={editingUser.name}
                    onChange={(e) => setEditingUser({...editingUser, name: e.target.value})}
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Yeni Şifre (boş bırakılırsa değişmez):</label>
                  <input
                    type="password"
                    value={editingUser.password}
                    onChange={(e) => setEditingUser({...editingUser, password: e.target.value})}
                  />
                </div>
                <div className="form-group">
                  <label>Rol:</label>
                  <select
                    value={editingUser.role}
                    onChange={(e) => setEditingUser({...editingUser, role: e.target.value})}
                  >
                    <option value="temizlikci">Temizlikçi</option>
                    <option value="admin">Admin</option>
                  </select>
                </div>
                <div className="modal-buttons">
                  <button type="submit" className="btn-primary">Güncelle</button>
                  <button 
                    type="button" 
                    className="btn-secondary" 
                    onClick={() => setEditingUser(null)}
                  >
                    İptal
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    </>
  );
};

export default AdminPanel;
