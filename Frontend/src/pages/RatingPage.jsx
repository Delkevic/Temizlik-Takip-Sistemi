import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import '../components/RatingPage.css';

const RatingPage = () => {
  const { toiletId } = useParams();
  const navigate = useNavigate();
  const [rating, setRating] = useState(0);
  const [hoveredRating, setHoveredRating] = useState(0);
  const [selectedProblems, setSelectedProblems] = useState([]);
  const [otherProblemText, setOtherProblemText] = useState('');

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
          // Ana sayfaya yönlendir
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

  return (
    <div className="rating-page">
      <div className="rating-container">
        <h1 className="page-title">Temizlik Değerlendirme - Tuvalet {toiletId}</h1>
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