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

// Toilet tuvaletleri temsil eden model
type Toilet struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Location  string    `json:"location" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToiletStatus tuvalet durumu için struct
type ToiletStatus struct {
	Toilet        Toilet        `json:"toilet"`
	LastRating    *Rating       `json:"last_rating,omitempty"`
	HasProblems   bool          `json:"has_problems"`
	ProblemCount  int           `json:"problem_count"`
	LastChecked   *time.Time    `json:"last_checked,omitempty"`
	CleaningTask  *CleaningTask `json:"cleaning_task,omitempty"`
	AverageRating float64       `json:"average_rating"`
	TotalRatings  int           `json:"total_ratings"`
}

// CleaningTask temizlik görevi için model
type CleaningTask struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	ToiletID    int        `json:"toilet_id" gorm:"not null"`
	CleanerID   uint       `json:"cleaner_id" gorm:"not null"`
	CleanerName string     `json:"cleaner_name" gorm:"not null;size:255"`
	Status      string     `json:"status" gorm:"not null;default:'assigned';size:50"` // assigned, in_progress, completed
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CleaningTaskRequest temizlik görevi isteği için struct
type CleaningTaskRequest struct {
	ToiletID int `json:"toilet_id" binding:"required,min=1"`
}

// CleaningTaskResponse temizlik görevi yanıtı için struct
type CleaningTaskResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Task    *CleaningTask `json:"task,omitempty"`
}

// CreateUserRequest yeni kullanıcı oluşturmak için struct
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

// UpdateUserRequest kullanıcı güncellemek için struct
type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}

// UserResponse kullanıcı response için struct
type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

// UsersResponse çoklu kullanıcı response için struct
type UsersResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Users   []User `json:"users,omitempty"`
}

// CleanerStats temizlikçi istatistikleri için struct
type CleanerStats struct {
	CleanerID           uint    `json:"cleaner_id"`
	CleanerName         string  `json:"cleaner_name"`
	TotalCompletedTasks int64   `json:"total_completed_tasks"`
	AverageCleaningTime float64 `json:"average_cleaning_time"` // dakika olarak
	LastWeekTasks       int64   `json:"last_week_tasks"`
	LastMonthTasks      int64   `json:"last_month_tasks"`
	IsActive            bool    `json:"is_active"`
	OngoingTasks        int64   `json:"ongoing_tasks"`
	FastestCleaningTime float64 `json:"fastest_cleaning_time"` // dakika olarak
	SlowestCleaningTime float64 `json:"slowest_cleaning_time"` // dakika olarak
	TotalCleaningTime   float64 `json:"total_cleaning_time"`   // dakika olarak
}

// SystemStats sistem geneli istatistikleri için struct
type SystemStats struct {
	TotalToilets        int64   `json:"total_toilets"`
	ActiveToilets       int64   `json:"active_toilets"`
	ToiletsWithProblems int     `json:"toilets_with_problems"`
	TotalCleaners       int64   `json:"total_cleaners"`
	ActiveCleaners      int64   `json:"active_cleaners"`
	TotalRatings        int64   `json:"total_ratings"`
	AverageRating       float64 `json:"average_rating"`
	CompletedTasksToday int64   `json:"completed_tasks_today"`
	OngoingTasks        int64   `json:"ongoing_tasks"`
}

// StatsResponse istatistik yanıtı için struct
type StatsResponse struct {
	Success      bool           `json:"success"`
	Message      string         `json:"message"`
	SystemStats  *SystemStats   `json:"system_stats,omitempty"`
	CleanerStats []CleanerStats `json:"cleaner_stats,omitempty"`
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

// RatingDetail değerlendirme detayları için struct
type RatingDetail struct {
	Rating    `json:",inline"`
	Problems  []string  `json:"problem_texts"`
	CreatedAt time.Time `json:"created_at"`
}

// PaginatedRatingsResponse sayfalı değerlendirme yanıtı için struct
type PaginatedRatingsResponse struct {
	Success     bool           `json:"success"`
	Message     string         `json:"message"`
	Data        []RatingDetail `json:"data"`
	ToiletID    int            `json:"toilet_id"`
	Page        int            `json:"page"`
	Limit       int            `json:"limit"`
	TotalCount  int            `json:"total_count"`
	TotalPages  int            `json:"total_pages"`
	HasNext     bool           `json:"has_next"`
	HasPrevious bool           `json:"has_previous"`
}
