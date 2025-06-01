package domain

import "github.com/google/uuid"

type AudioPhrase struct {
	ID          uuid.UUID `json:"id"`
	PathToAudio string    `json:"path_to_audio"`
	PhraseID    uuid.UUID `json:"phrase_id"`
	Accent      string    `json:"accent"`
	Noise       float64   `json:"noise"`
}
