package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Director    string         `json:"director"`
	Year        int            `json:"year"`
	Genre       string         `json:"genre"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Reviews     []Review       `json:"reviews,omitempty"`
}

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MovieID   uint           `json:"movie_id" gorm:"not null;index"`
	Rating    int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 10"`
	Comment   string         `json:"comment"`
	WatchedAt time.Time      `json:"watched_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Movie     Movie          `json:"movie,omitempty"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Movie{}, &Review{})
}
