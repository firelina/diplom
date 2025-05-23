package domain

import "github.com/google/uuid"

type PhraseType struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
