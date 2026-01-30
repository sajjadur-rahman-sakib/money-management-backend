package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	BookID      uuid.UUID `json:"book_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // "cash_in" or "cash_out"
	Amount      float64   `json:"amount" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Book        Book      `json:"book" gorm:"foreignKey:BookID"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
