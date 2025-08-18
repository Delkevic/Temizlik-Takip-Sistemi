import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import '../components/RatingPage.css';

const RatingPage = () => {
  const { toiletId: paramToiletId } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const [toiletId, setToiletId] = useState(null);
  const [toiletName, setToiletName] = useState('');
  const [toilets, setToilets] = useState([]);
  const [showToiletSelection, setShowToiletSelection] = useState(false);
  const [rating, setRating] = useState(0);
  const [hoveredRating, setHoveredRating] = useState(0);
  const [selectedProblems, setSelectedProblems] = useState([]);
  const [otherProblemText, setOtherProblemText] = useState('');

  useEffect(() => {
    // URL parametresinden veya query parametresinden toilet ID'sini al
    const searchParams = new URLSearchParams(location.search);
    const queryToiletId = searchParams.get('toilet');
    
    const finalToiletId = paramToiletId || queryToiletId;
    
    if (finalToiletId) {
      setToiletId(finalToiletId);
      fetchToiletInfo(finalToiletId);
    } else {
      // Tuvalet seçilmemişse, seçim ekranını göster
      setShowToiletSelection(true);
      fetchToilets();
    }
  }, [paramToiletId, location.search]);

  const fetchToilets = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/toilets');
      const data = await response.json();
      
      if (data.success) {
        setToilets(data.data);
      }
    } catch (error) {
      console.error('Tuvaletler getirilemedi:', error);
    }
  };

  const fetchToiletInfo = async (id) => {
    try {
      const response = await fetch('http://localhost:8080/api/toilets');
      const data = await response.json();
      
      if (data.success) {
        const toilet = data.data.find(t => t.id === parseInt(id));
        if (toilet) {
          setToiletName(toilet.name);
        }
      }
    } catch (error) {
      console.error('Tuvalet bilgisi getirilemedi:', error);
    }
  };

  const selectToilet = (selectedToiletId, selectedToiletName) => {
    setToiletId(selectedToiletId);
    setToiletName(selectedToiletName);
    setShowToiletSelection(false);
  };

  const handleRatingClick = (value) => {
    setRating(value);
  };

  const handleMouseEnter = (value) => {
    setHoveredRating(value);
  };

  const handleMouseLeave = () => {
    setHoveredRating(0);
  };

  const handleProblemToggle = (problemId) => {
    setSelectedProblems(prev => 
      prev.includes(problemId) 
        ? prev.filter(id => id !== problemId)
        : [...prev, problemId]
    );
    
    // Eğer "Diğer" seçeneği kaldırılıyorsa, text'i de temizle
    if (problemId === 6 && selectedProblems.includes(problemId)) {
      setOtherProblemText('');
    }
  };

  const handleSubmit = async () => {
    if (rating > 0) {
      try {
        const response = await fetch('http://localhost:8080/api/rating', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            toilet_id: parseInt(toiletId),
            rating: rating,
            problems: selectedProblems,
            other_text: otherProblemText
          })
        });

        const data = await response.json();

        if (data.success) {
          alert('Değerlendirmeniz başarıyla kaydedildi!');
          // Formu sıfırla
          setRating(0);
          setSelectedProblems([]);
          setOtherProblemText('');
          // Temizlikçi paneline yönlendir
          navigate('/');
        } else {
          alert('Hata: ' + data.message);
        }
      } catch (error) {
        console.error('API Hatası:', error);
        alert('Bağlantı hatası oluştu. Lütfen tekrar deneyin.');
      }
    } else {
      alert('Lütfen bir puan seçin!');
    }
  };

  if (showToiletSelection) {
    return (
      <div className="rating-page">
        <div className="rating-container">
          <h1 className="page-title">Tuvalet Seçimi</h1>
          <h2 className="page-subtitle">Değerlendirmek istediğiniz tuvaleti seçin</h2>
          
          <div className="toilet-selection-grid">
            {toilets.map((toilet) => (
              <button
                key={toilet.id}
                className="toilet-selection-button"
                onClick={() => selectToilet(toilet.id, toilet.name)}
              >
                <h3>{toilet.name}</h3>
                <p>📍 {toilet.location}</p>
              </button>
            ))}
          </div>
          
          <button 
            className="back-button"
            onClick={() => navigate('/')}
          >
            Geri Dön
          </button>
        </div>
      </div>
    );
  }

  if (!toiletId) {
    return (
      <div className="rating-page">
        <div className="rating-container">
          <h1>Geçersiz tuvalet seçimi</h1>
          <button onClick={() => navigate('/')} className="submit-button">
            Geri Dön
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="rating-page">
      <div className="rating-container">
        <h1 className="page-title">
          Temizlik Değerlendirme - {toiletName || `Tuvalet ${toiletId}`}
        </h1>
        <h2 className="page-subtitle">Lütfen temizlik hizmetini değerlendirin</h2>
        
        <div className="star-rating">
          {[1, 2, 3, 4, 5].map((star) => (
            <span
              key={star}
              className={`star ${
                star <= (hoveredRating || rating) ? 'star-filled' : 'star-empty'
              }`}
              onClick={() => handleRatingClick(star)}
              onMouseEnter={() => handleMouseEnter(star)}
              onMouseLeave={handleMouseLeave}
            >
              ★
            </span>
          ))}
        </div>

        {rating > 0 && (
          <p className="rating-text">
            Seçtiğiniz puan: {rating} / 5
          </p>
        )}

        <h2 className="page-subtitle">Sorunlardan biri varsa seçiniz yoksa devam edin</h2>

        <div className="problem-options">
          {[1, 2, 3, 4, 5, 6].map((problemId) => (
            <button
              key={problemId}
              className={`problem-button ${
                selectedProblems.includes(problemId) ? 'problem-selected' : 'problem-unselected'
              }`}
              onClick={() => handleProblemToggle(problemId)}
            >
              {problemId === 1 ? "Tuvalet Kağıdı yok" : ""}
              {problemId === 2 ? "Sabun yok" : ""}
              {problemId === 3 ? "Peçete yok" : ""}
              {problemId === 4 ? "Çöp kutusu dolu" : ""}
              {problemId === 5 ? "Klozet kirli" : ""}
              {problemId === 6 ? "Diğer" : ""}
            </button>
          ))}
        </div>

        {selectedProblems.includes(6) && (
          <div className="other-problem-input">
            <label htmlFor="otherProblem" className="input-label">
              Lütfen sorunu detaylandırın:
            </label>
            <textarea
              id="otherProblem"
              className="other-problem-textarea"
              value={otherProblemText}
              onChange={(e) => setOtherProblemText(e.target.value)}
              placeholder="Karşılaştığınız sorunu buraya yazın..."
              rows={3}
            />
          </div>
        )}

        <button 
          className="submit-button"
          onClick={handleSubmit}
          disabled={rating === 0}
        >
          Gönder
        </button>
      </div>
    </div>
  );
};

export default RatingPage;