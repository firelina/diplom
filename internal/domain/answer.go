package domain

import "github.com/google/uuid"

type Answer struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AudioAnswerID uuid.UUID `json:"audio_answer_id"`
	Text          string    `json:"text"`
	IsCorrect     bool      `json:"is_correct"`
}
