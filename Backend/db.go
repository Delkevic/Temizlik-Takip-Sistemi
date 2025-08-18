package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase MySQL veritabanı bağlantısını başlatır
func InitDatabase() {
	// .env dosyasını yükle
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenemedi: ", err)
	}

	// Environment variables'dan bağlantı bilgilerini al
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// MySQL bağlantı string'i oluştur
	// Format: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci&interpolateParams=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Veritabanı bağlantısı başarısız: ", err)
	}

	// Veritabanı charset'ini ayarla
	charsetCommands := []string{
		"SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci",
		"SET character_set_connection=utf8mb4",
		"SET character_set_results=utf8mb4",
		"SET character_set_client=utf8mb4",
		"SET character_set_server=utf8mb4",
		"SET collation_connection=utf8mb4_unicode_ci",
		"SET collation_server=utf8mb4_unicode_ci",
	}

	for _, cmd := range charsetCommands {
		if err := DB.Exec(cmd).Error; err != nil {
			log.Printf("Charset ayarı hatası (%s): %v", cmd, err)
		}
	}

	// Veritabanı tablolarını otomatik oluştur
	err = DB.AutoMigrate(&User{}, &Rating{}, &Toilet{})
	if err != nil {
		log.Fatal("Tablo oluşturma hatası: ", err)
	}

	// CleaningTask tablosunu manuel oluştur (charset sorunu için)
	recreateCleaningTaskTable()

	// İlk tuvaletleri oluştur (sadece bir kez)
	createInitialToilets()

	log.Println("Veritabanı bağlantısı başarılı!")
}

// GetDB veritabanı bağlantısını döndürür
func GetDB() *gorm.DB {
	return DB
}

// createInitialToilets ilk 6 tuvaleti oluşturur
func createInitialToilets() {
	// Tuvalet sayısını kontrol et
	var count int64
	DB.Model(&Toilet{}).Count(&count)

	// Eğer hiç tuvalet yoksa, 6 tane oluştur
	if count == 0 {
		toilets := []Toilet{
			{ID: 1, Name: "1. Kat Erkek Tuvaleti", Location: "1. Kat", IsActive: true},
			{ID: 2, Name: "1. Kat Kadın Tuvaleti", Location: "1. Kat", IsActive: true},
			{ID: 3, Name: "2. Kat Erkek Tuvaleti", Location: "2. Kat", IsActive: true},
			{ID: 4, Name: "2. Kat Kadın Tuvaleti", Location: "2. Kat", IsActive: true},
			{ID: 5, Name: "3. Kat Erkek Tuvaleti", Location: "3. Kat", IsActive: true},
			{ID: 6, Name: "3. Kat Kadın Tuvaleti", Location: "3. Kat", IsActive: true},
		}

		for _, toilet := range toilets {
			if err := DB.Create(&toilet).Error; err != nil {
				log.Printf("Tuvalet oluşturma hatası: %v", err)
			}
		}

		log.Println("İlk 6 tuvalet başarıyla oluşturuldu!")
	}
}

// recreateCleaningTaskTable CleaningTask tablosunu doğru charset ile yeniden oluşturur
func recreateCleaningTaskTable() {
	// Önce mevcut tabloyu kontrol et
	var tableExists bool
	DB.Raw("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'cleaning_tasks'").Scan(&tableExists)

	if tableExists {
		log.Println("Mevcut cleaning_tasks tablosu siliniyor...")
		DB.Exec("DROP TABLE cleaning_tasks")
	}

	// Tabloyu doğru charset ve collation ile oluştur
	createTableSQL := `
	CREATE TABLE cleaning_tasks (
		id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		toilet_id INT NOT NULL,
		cleaner_id INT UNSIGNED NOT NULL,
		cleaner_name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
		status VARCHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'assigned',
		started_at TIMESTAMP NULL,
		completed_at TIMESTAMP NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_toilet_status (toilet_id, status)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

	result := DB.Exec(createTableSQL)
	if result.Error != nil {
		log.Printf("CleaningTask tablosu oluşturma hatası: %v", result.Error)
	} else {
		log.Println("CleaningTask tablosu başarıyla UTF8MB4 ile oluşturuldu!")

		// Test verisi ekle
		testSQL := "INSERT INTO cleaning_tasks (toilet_id, cleaner_id, cleaner_name, status) VALUES (999, 999, 'Test Görevlisi Örnekleme', 'test')"
		testResult := DB.Exec(testSQL)
		if testResult.Error != nil {
			log.Printf("UTF8 test hatası: %v", testResult.Error)
		} else {
			log.Println("UTF8 test başarılı - Türkçe karakterler destekleniyor!")
			// Test verisini sil
			DB.Exec("DELETE FROM cleaning_tasks WHERE toilet_id = 999")
		}
	}
}
