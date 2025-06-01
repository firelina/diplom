package models

type CreatePhraseStreamRequest struct {
	PhraseID   string  `json:"phrase_id"`
	Path       string  `json:"path"`
	ScenarioID string  `json:"scenario_id"`
	Accent     string  `json:"accent"`
	Noise      float64 `json:"noise"`
}
