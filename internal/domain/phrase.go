package domain

import "github.com/google/uuid"

type Phrase struct {
	ID     uuid.UUID `json:"id"`
	Text   string    `json:"text"`
	TypeID uuid.UUID `json:"type_id"`
}
