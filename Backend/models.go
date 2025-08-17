package main

import (
	"time"
)

// User çalışanları temsil eden model
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // JSON'da gizli
	Name      string    `json:"name" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;default:'temizlikci'"` // admin, temizlikci
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRequest giriş isteği için struct
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse giriş yanıtı için struct
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// Rating veritabanında puanlamaları temsil eden model
type Rating struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ToiletID  int       `json:"toilet_id" gorm:"not null"`
	Rating    int       `json:"rating" gorm:"not null"`
	Problems  string    `json:"problems"`   // JSON string olarak problemi IDs saklanacak
	OtherText string    `json:"other_text"` // "Diğer" seçeneği için metin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RatingRequest frontend'den gelen rating verisini temsil eden struct
type RatingRequest struct {
	ToiletID  int    `json:"toilet_id" binding:"required,min=1"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Problems  []int  `json:"problems"`
	OtherText string `json:"other_text"`
}

// RatingResponse API'den dönen response yapısı
type RatingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      uint   `json:"id,omitempty"`
}

// Problem türlerini tanımlayan sabitler
var ProblemTypes = map[int]string{
	1: "Tuvalet Kağıdı yok",
	2: "Sabun yok",
	3: "Peçete yok",
	4: "Çöp kutusu dolu",
	5: "Klozet kirli",
	6: "Diğer",
}
