package models

type CreateScenarioRequest struct {
	Title  string `json:"title"`
	UserID string `json:"user_id"`
}
