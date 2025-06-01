package models

import "github.com/google/uuid"

type CreateAnswerResponse struct {
	AnswerID  uuid.UUID `json:"answer_id"`
	IsCorrect bool      `json:"is_correct"`
	Text      string    `json:"text"`
}
