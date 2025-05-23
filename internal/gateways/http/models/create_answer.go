package models

type CreateAnswerRequest struct {
	Path           string `json:"path"`
	UserID         string `json:"user_id"`
	PhraseStreamID string `json:"phrase_stream_id"`
}
