package domain

import "github.com/google/uuid"

type PhraseStream struct {
	ID            uuid.UUID `json:"id"`
	AudioPhraseID uuid.UUID `json:"audio_phrase_id"`
	ScenarioID    uuid.UUID `json:"scenario_id"`
	AnswerID      uuid.UUID `json:"answer_id"`
	PhraseID      uuid.UUID `json:"phrase_id"`
	Status        string    `json:"status"`
}
