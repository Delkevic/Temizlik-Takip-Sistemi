package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Veritabanı bağlantısını başlat
	InitDatabase()

	// Gin router'ı oluştur
	router := gin.Default()

	// Route'ları ayarla
	SetupRoutes(router)

	// Sunucuyu 8080 portunda başlat
	log.Println("Server 8080 portunda başlatılıyor...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server başlatılamadı: ", err)
	}
}
