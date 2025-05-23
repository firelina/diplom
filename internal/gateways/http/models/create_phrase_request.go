package models

import "github.com/google/uuid"

type CreatePhraseRequest struct {
	Text   string    `json:"text"`
	TypeID uuid.UUID `json:"type_id"`
}
