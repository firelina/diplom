package domain

import (
	"github.com/google/uuid"
	"time"
)

type Scenario struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	UserID    uuid.UUID  `json:"user_id"`
}
