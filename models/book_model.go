package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	Balance   float64   `json:"balance" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	//Transactions []Transaction `json:"transactions" gorm:"foreignKey:BookID"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}
