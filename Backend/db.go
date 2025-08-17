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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanı bağlantısı başarısız: ", err)
	}

	// Veritabanı tablolarını otomatik oluştur
	err = DB.AutoMigrate(&User{}, &Rating{})
	if err != nil {
		log.Fatal("Tablo oluşturma hatası: ", err)
	}

	log.Println("Veritabanı bağlantısı başarılı!")
}

// GetDB veritabanı bağlantısını döndürür
func GetDB() *gorm.DB {
	return DB
}
