package models

type PhraseProgress struct {
	Phrase             string `json:"phrase"`
	PhraseStreamStatus string `json:"phrase_stream_status"`
	ScenarioStatus     string `json:"scenario_status"`
}

type Progress []*PhraseProgress
