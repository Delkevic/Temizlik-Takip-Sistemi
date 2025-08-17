package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupRoutes API route'larını ayarlar
func SetupRoutes(router *gin.Engine) {
	// CORS middleware ekle
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/login", login)

		// Rating routes
		api.POST("/rating", createRating)
		api.GET("/ratings", getRatings)
		api.GET("/rating/:id", getRating)
		api.GET("/toilet/:toiletId/ratings", getToiletRatings)
	}
}

// hashPassword şifreyi SHA256 ile hashler
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// login kullanıcı girişi yapar
func login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "Geçersiz veri formatı: " + err.Error(),
		})
		return
	}

	// Kullanıcıyı veritabanında ara
	var user User
	if err := DB.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Kullanıcı adı veya şifre hatalı",
		})
		return
	}

	// Şifreyi kontrol et (basit SHA256 hash)
	hashedPassword := hashPassword(req.Password)
	if user.Password != hashedPassword {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Kullanıcı adı veya şifre hatalı",
		})
		return
	}

	// Basit token oluştur (production'da JWT kullanın)
	token := hashPassword(user.Username + "salt123")

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Giriş başarılı",
		Token:   token,
		User:    &user,
	})
}

// createRating yeni bir puanlama oluşturur
func createRating(c *gin.Context) {
	var req RatingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, RatingResponse{
			Success: false,
			Message: "Geçersiz veri formatı: " + err.Error(),
		})
		return
	}

	// Problems slice'ını JSON string'e çevir
	problemsJSON, err := json.Marshal(req.Problems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RatingResponse{
			Success: false,
			Message: "Sorun verileri işlenirken hata oluştu",
		})
		return
	}

	// Yeni rating oluştur
	rating := Rating{
		ToiletID:  req.ToiletID,
		Rating:    req.Rating,
		Problems:  string(problemsJSON),
		OtherText: req.OtherText,
	}

	// Veritabanına kaydet
	if err := DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, RatingResponse{
			Success: false,
			Message: "Veritabanı kayıt hatası: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, RatingResponse{
		Success: true,
		Message: "Puanlama başarıyla kaydedildi",
		ID:      rating.ID,
	})
}

// getRatings tüm puanlamaları getirir
func getRatings(c *gin.Context) {
	var ratings []Rating

	if err := DB.Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Veriler getirilirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ratings,
	})
}

// getRating belirli bir puanlamayı getirir
func getRating(c *gin.Context) {
	id := c.Param("id")

	ratingID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Geçersiz ID formatı",
		})
		return
	}

	var rating Rating
	if err := DB.First(&rating, uint(ratingID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Puanlama bulunamadı",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rating,
	})
}

// getToiletRatings belirli bir tuvalete ait tüm puanlamaları getirir
func getToiletRatings(c *gin.Context) {
	toiletIdStr := c.Param("toiletId")

	toiletID, err := strconv.Atoi(toiletIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Geçersiz tuvalet ID formatı",
		})
		return
	}

	var ratings []Rating
	if err := DB.Where("toilet_id = ?", toiletID).Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Veriler getirilirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      ratings,
		"toilet_id": toiletID,
		"count":     len(ratings),
	})
}
