package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes API route'larını ayarlar
func SetupRoutes(router *gin.Engine) {
	// CORS middleware ekle
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-ID, X-User-Name")

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

		// Toilet routes
		api.GET("/toilets", getToilets)
		api.GET("/toilets/status", getToiletsStatus)

		// Cleaning task routes
		api.POST("/cleaning/start", startCleaningTask)
		api.PUT("/cleaning/begin/:id", beginCleaningTask)
		api.PUT("/cleaning/complete/:id", completeCleaningTask)
		api.GET("/cleaning/tasks", getCleaningTasks)

		// Admin routes - User management
		api.GET("/admin/users", getUsers)
		api.POST("/admin/users", createUser)
		api.PUT("/admin/users/:id", updateUser)
		api.DELETE("/admin/users/:id", deleteUser)

		// Admin routes - Statistics
		api.GET("/admin/stats", getAdminStats)
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

// getToilets tüm aktif tuvaletleri getirir
func getToilets(c *gin.Context) {
	var toilets []Toilet

	if err := DB.Where("is_active = ?", true).Find(&toilets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Veriler getirilirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toilets,
	})
}

// getToiletsStatus tüm tuvaletlerin durumunu getirir
func getToiletsStatus(c *gin.Context) {
	var toilets []Toilet

	// Aktif tuvaletleri getir
	if err := DB.Where("is_active = ?", true).Find(&toilets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Tuvaletler getirilirken hata oluştu: " + err.Error(),
		})
		return
	}

	var toiletStatuses []ToiletStatus

	for _, toilet := range toilets {
		var lastRating Rating
		var hasLastRating bool

		// Son puanlamayı getir
		if err := DB.Where("toilet_id = ?", toilet.ID).Order("created_at DESC").First(&lastRating).Error; err == nil {
			hasLastRating = true
		}

		// Ortalama puanı hesapla
		var avgRating struct {
			Average float64
			Count   int64
		}
		DB.Raw("SELECT AVG(rating) as average, COUNT(*) as count FROM ratings WHERE toilet_id = ?", toilet.ID).Scan(&avgRating)

		// Problem sayısını hesapla
		var problemCount int
		hasProblems := false
		var lastChecked *time.Time

		if hasLastRating {
			lastChecked = &lastRating.CreatedAt

			// Problems JSON string'ini parse et
			var problems []int
			if err := json.Unmarshal([]byte(lastRating.Problems), &problems); err == nil {
				problemCount = len(problems)
				hasProblems = problemCount > 0
			}
		}

		// Aktif temizlik görevini kontrol et
		var cleaningTask CleaningTask
		hasActiveTask := false
		if err := DB.Where("toilet_id = ? AND status IN ?", toilet.ID, []string{"assigned", "in_progress"}).First(&cleaningTask).Error; err == nil {
			hasActiveTask = true
		}

		status := ToiletStatus{
			Toilet:        toilet,
			HasProblems:   hasProblems,
			ProblemCount:  problemCount,
			LastChecked:   lastChecked,
			AverageRating: avgRating.Average,
			TotalRatings:  int(avgRating.Count),
		}

		if hasLastRating {
			status.LastRating = &lastRating
		}

		if hasActiveTask {
			status.CleaningTask = &cleaningTask
		}

		toiletStatuses = append(toiletStatuses, status)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    toiletStatuses,
	})
}

// startCleaningTask temizlik görevini başlatır
func startCleaningTask(c *gin.Context) {
	var req CleaningTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CleaningTaskResponse{
			Success: false,
			Message: "Geçersiz veri formatı: " + err.Error(),
		})
		return
	}

	// Bu tuvalet için aktif bir temizlik görevi var mı kontrol et
	var existingTask CleaningTask
	if err := DB.Where("toilet_id = ? AND status IN ?", req.ToiletID, []string{"assigned", "in_progress"}).First(&existingTask).Error; err == nil {
		c.JSON(http.StatusConflict, CleaningTaskResponse{
			Success: false,
			Message: "Bu tuvalet için zaten aktif bir temizlik görevi var",
		})
		return
	}

	// Authorization header'ından kullanıcı bilgisini al
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, CleaningTaskResponse{
			Success: false,
			Message: "Yetkilendirme başlığı eksik",
		})
		return
	}

	// Basit token doğrulama (gerçek uygulamada JWT kullanılmalı)
	// Şimdilik header'dan user_id ve user_name alıyoruz
	cleanerIDHeader := c.GetHeader("X-User-ID")
	cleanerNameHeader := c.GetHeader("X-User-Name")

	var cleanerID uint = 1
	var cleanerName string = "Temizlik Görevlisi"

	if cleanerIDHeader != "" {
		if id, err := strconv.ParseUint(cleanerIDHeader, 10, 32); err == nil {
			cleanerID = uint(id)
		}
	}

	if cleanerNameHeader != "" {
		cleanerName = cleanerNameHeader
	}

	// Yeni temizlik görevi oluştur
	task := CleaningTask{
		ToiletID:    req.ToiletID,
		CleanerID:   cleanerID,
		CleanerName: cleanerName,
		Status:      "assigned",
		StartedAt:   nil,
		CompletedAt: nil,
	}

	if err := DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, CleaningTaskResponse{
			Success: false,
			Message: "Temizlik görevi oluşturulurken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, CleaningTaskResponse{
		Success: true,
		Message: "Temizlik görevi başarıyla oluşturuldu",
		Task:    &task,
	})
}

// beginCleaningTask temizlik görevini başlat durumuna getirir
func beginCleaningTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, CleaningTaskResponse{
			Success: false,
			Message: "Geçersiz görev ID formatı",
		})
		return
	}

	var task CleaningTask
	if err := DB.First(&task, uint(taskID)).Error; err != nil {
		c.JSON(http.StatusNotFound, CleaningTaskResponse{
			Success: false,
			Message: "Temizlik görevi bulunamadı",
		})
		return
	}

	// Görevin durumunu güncelle
	now := time.Now()
	task.Status = "in_progress"
	task.StartedAt = &now

	if err := DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, CleaningTaskResponse{
			Success: false,
			Message: "Temizlik görevi güncellenirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CleaningTaskResponse{
		Success: true,
		Message: "Temizlik görevi başlatıldı",
		Task:    &task,
	})
}

// completeCleaningTask temizlik görevini tamamlar
func completeCleaningTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, CleaningTaskResponse{
			Success: false,
			Message: "Geçersiz görev ID formatı",
		})
		return
	}

	var task CleaningTask
	if err := DB.First(&task, uint(taskID)).Error; err != nil {
		c.JSON(http.StatusNotFound, CleaningTaskResponse{
			Success: false,
			Message: "Temizlik görevi bulunamadı",
		})
		return
	}

	// Transaction ile işlemi gerçekleştir
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Görevin durumunu güncelle
	now := time.Now()
	task.Status = "completed"
	task.CompletedAt = &now

	if task.StartedAt == nil {
		task.StartedAt = &now
	}

	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, CleaningTaskResponse{
			Success: false,
			Message: "Temizlik görevi güncellenirken hata oluştu: " + err.Error(),
		})
		return
	}

	// Temizlik tamamlandığında otomatik olarak temiz bir rating oluştur
	// Böylece tuvalet "temiz" olarak gözükecek
	cleanRating := Rating{
		ToiletID:  task.ToiletID,
		Rating:    5,    // Temizlik tamamlandığında maksimum puan ver
		Problems:  "[]", // Boş problem listesi (JSON array)
		OtherText: "",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tx.Create(&cleanRating).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, CleaningTaskResponse{
			Success: false,
			Message: "Temizlik sonrası değerlendirme oluşturulurken hata oluştu: " + err.Error(),
		})
		return
	}

	// Transaction'ı tamamla
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, CleaningTaskResponse{
			Success: false,
			Message: "İşlem tamamlanırken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CleaningTaskResponse{
		Success: true,
		Message: "Temizlik görevi başarıyla tamamlandı",
		Task:    &task,
	})
}

// getCleaningTasks temizlik görevlerini getirir
func getCleaningTasks(c *gin.Context) {
	status := c.Query("status")
	toiletID := c.Query("toilet_id")

	query := DB.Model(&CleaningTask{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if toiletID != "" {
		query = query.Where("toilet_id = ?", toiletID)
	}

	var tasks []CleaningTask
	if err := query.Order("created_at DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Temizlik görevleri getirilirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tasks,
	})
}

// getUsers tüm kullanıcıları getirir (sadece admin erişimi)
func getUsers(c *gin.Context) {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UsersResponse{
			Success: false,
			Message: "Kullanıcılar getirilirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UsersResponse{
		Success: true,
		Message: "Kullanıcılar başarıyla getirildi",
		Users:   users,
	})
}

// createUser yeni kullanıcı oluşturur (sadece admin erişimi)
func createUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "Geçersiz veri formatı: " + err.Error(),
		})
		return
	}

	// Kullanıcı adı kontrolü
	var existingUser User
	if err := DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, UserResponse{
			Success: false,
			Message: "Bu kullanıcı adı zaten kullanılıyor",
		})
		return
	}

	// Role varsayılan değeri
	if req.Role == "" {
		req.Role = "temizlikci"
	}

	// Şifreyi hashle
	hashedPassword := hashPassword(req.Password)

	// Yeni kullanıcı oluştur
	user := User{
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
		Role:     req.Role,
		IsActive: true,
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Success: false,
			Message: "Kullanıcı oluşturulurken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		Success: true,
		Message: "Kullanıcı başarıyla oluşturuldu",
		User:    &user,
	})
}

// updateUser kullanıcı bilgilerini günceller (sadece admin erişimi)
func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "Geçersiz kullanıcı ID formatı",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "Geçersiz veri formatı: " + err.Error(),
		})
		return
	}

	var user User
	if err := DB.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, UserResponse{
			Success: false,
			Message: "Kullanıcı bulunamadı",
		})
		return
	}

	// Kullanıcı adı değişikliği kontrolü
	if req.Username != "" && req.Username != user.Username {
		var existingUser User
		if err := DB.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, UserResponse{
				Success: false,
				Message: "Bu kullanıcı adı zaten kullanılıyor",
			})
			return
		}
		user.Username = req.Username
	}

	// Diğer alanları güncelle
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		user.Password = hashPassword(req.Password)
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Success: false,
			Message: "Kullanıcı güncellenirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		Message: "Kullanıcı başarıyla güncellendi",
		User:    &user,
	})
}

// deleteUser kullanıcıyı siler (sadece admin erişimi)
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Success: false,
			Message: "Geçersiz kullanıcı ID formatı",
		})
		return
	}

	var user User
	if err := DB.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, UserResponse{
			Success: false,
			Message: "Kullanıcı bulunamadı",
		})
		return
	}

	// Kullanıcıyı sil
	if err := DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserResponse{
			Success: false,
			Message: "Kullanıcı silinirken hata oluştu: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Success: true,
		Message: "Kullanıcı başarıyla silindi",
	})
}

// getAdminStats admin paneli için istatistikleri getirir
func getAdminStats(c *gin.Context) {
	// Sistem geneli istatistikler
	var systemStats SystemStats

	// Toplam tuvalet sayısı
	DB.Model(&Toilet{}).Count(&systemStats.TotalToilets)

	// Aktif tuvalet sayısı
	DB.Model(&Toilet{}).Where("is_active = ?", true).Count(&systemStats.ActiveToilets)

	// Problem olan tuvalet sayısı (son 24 saat içinde problem bildirilen)
	var problemToilets []int
	DB.Model(&Rating{}).
		Select("DISTINCT toilet_id").
		Where("problems IS NOT NULL AND problems != '' AND problems != '[]' AND created_at > ?", time.Now().Add(-24*time.Hour)).
		Pluck("toilet_id", &problemToilets)
	systemStats.ToiletsWithProblems = len(problemToilets)

	// Toplam temizlikçi sayısı
	DB.Model(&User{}).Where("role = ?", "temizlikci").Count(&systemStats.TotalCleaners)

	// Aktif temizlikçi sayısı
	DB.Model(&User{}).Where("role = ? AND is_active = ?", "temizlikci", true).Count(&systemStats.ActiveCleaners)

	// Toplam değerlendirme sayısı
	DB.Model(&Rating{}).Count(&systemStats.TotalRatings)

	// Ortalama puan
	var avgRating sql.NullFloat64
	DB.Model(&Rating{}).Select("AVG(rating)").Scan(&avgRating)
	if avgRating.Valid {
		systemStats.AverageRating = avgRating.Float64
	}

	// Bugün tamamlanan görev sayısı
	today := time.Now().Truncate(24 * time.Hour)
	DB.Model(&CleaningTask{}).
		Where("status = ? AND completed_at >= ?", "completed", today).
		Count(&systemStats.CompletedTasksToday)

	// Devam eden görev sayısı
	DB.Model(&CleaningTask{}).
		Where("status IN (?)", []string{"assigned", "in_progress"}).
		Count(&systemStats.OngoingTasks)

	// Temizlikçi istatistikleri
	var cleanerStats []CleanerStats

	// Tüm temizlikçileri al
	var cleaners []User
	DB.Where("role = ?", "temizlikci").Find(&cleaners)

	for _, cleaner := range cleaners {
		var stats CleanerStats
		stats.CleanerID = cleaner.ID
		stats.CleanerName = cleaner.Name
		stats.IsActive = cleaner.IsActive

		// Tamamlanan görev sayısı
		DB.Model(&CleaningTask{}).
			Where("cleaner_id = ? AND status = ?", cleaner.ID, "completed").
			Count(&stats.TotalCompletedTasks)

		// Ortalama temizlik süresi (dakika)
		var avgTime sql.NullFloat64
		DB.Model(&CleaningTask{}).
			Select("AVG(TIMESTAMPDIFF(MINUTE, started_at, completed_at))").
			Where("cleaner_id = ? AND status = ? AND started_at IS NOT NULL AND completed_at IS NOT NULL", cleaner.ID, "completed").
			Scan(&avgTime)
		if avgTime.Valid {
			stats.AverageCleaningTime = avgTime.Float64
		} else {
			stats.AverageCleaningTime = -1 // Hiç temizlik yapılmadığını belirtmek için
		}

		// En hızlı temizlik süresi
		var fastestTime sql.NullFloat64
		DB.Model(&CleaningTask{}).
			Select("MIN(TIMESTAMPDIFF(MINUTE, started_at, completed_at))").
			Where("cleaner_id = ? AND status = ? AND started_at IS NOT NULL AND completed_at IS NOT NULL", cleaner.ID, "completed").
			Scan(&fastestTime)
		if fastestTime.Valid {
			stats.FastestCleaningTime = fastestTime.Float64
		} else {
			stats.FastestCleaningTime = -1 // Hiç temizlik yapılmadığını belirtmek için
		}

		// En yavaş temizlik süresi
		var slowestTime sql.NullFloat64
		DB.Model(&CleaningTask{}).
			Select("MAX(TIMESTAMPDIFF(MINUTE, started_at, completed_at))").
			Where("cleaner_id = ? AND status = ? AND started_at IS NOT NULL AND completed_at IS NOT NULL", cleaner.ID, "completed").
			Scan(&slowestTime)
		if slowestTime.Valid {
			stats.SlowestCleaningTime = slowestTime.Float64
		} else {
			stats.SlowestCleaningTime = -1 // Hiç temizlik yapılmadığını belirtmek için
		}

		// Toplam temizlik süresi
		var totalTime sql.NullFloat64
		DB.Model(&CleaningTask{}).
			Select("SUM(TIMESTAMPDIFF(MINUTE, started_at, completed_at))").
			Where("cleaner_id = ? AND status = ? AND started_at IS NOT NULL AND completed_at IS NOT NULL", cleaner.ID, "completed").
			Scan(&totalTime)
		if totalTime.Valid {
			stats.TotalCleaningTime = totalTime.Float64
		}

		// Son hafta görev sayısı
		lastWeek := time.Now().Add(-7 * 24 * time.Hour)
		DB.Model(&CleaningTask{}).
			Where("cleaner_id = ? AND status = ? AND completed_at >= ?", cleaner.ID, "completed", lastWeek).
			Count(&stats.LastWeekTasks)

		// Son ay görev sayısı
		lastMonth := time.Now().Add(-30 * 24 * time.Hour)
		DB.Model(&CleaningTask{}).
			Where("cleaner_id = ? AND status = ? AND completed_at >= ?", cleaner.ID, "completed", lastMonth).
			Count(&stats.LastMonthTasks)

		// Devam eden görev sayısı
		DB.Model(&CleaningTask{}).
			Where("cleaner_id = ? AND status IN (?)", cleaner.ID, []string{"assigned", "in_progress"}).
			Count(&stats.OngoingTasks)

		cleanerStats = append(cleanerStats, stats)
	}

	c.JSON(http.StatusOK, StatsResponse{
		Success:      true,
		Message:      "İstatistikler başarıyla getirildi",
		SystemStats:  &systemStats,
		CleanerStats: cleanerStats,
	})
}
