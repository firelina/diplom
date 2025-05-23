package domain

import (
	"github.com/google/uuid"
	"time"
)

type AudioAnswer struct {
	ID          uuid.UUID `json:"id"`
	PathToAudio string    `json:"path_to_audio"`
	RecordTime  time.Time `json:"record_time"`
}
