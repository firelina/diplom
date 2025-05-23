package handlers

import (
	"diplom/internal/domain"
	"diplom/internal/gateways/http/models"
	"diplom/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type PhraseStreamHandler struct {
	phraseStreamService *services.PhraseStreamService
}

func NewPhraseStreamHandler(s *services.PhraseStreamService) *PhraseStreamHandler {
	return &PhraseStreamHandler{phraseStreamService: s}
}

// CreatePhraseStream godoc
// @Summary      Create a phrase stream
// @Description  Initializes a new phrase stream and stores associated audio
// @Tags         scenarios
// @Accept       json
// @Produce      json
// @Param        phraseStream  body      models.CreatePhraseStreamRequest  true  "Phrase stream data"
// @Success      201           {object}  string                             "Created phrase stream ID"
// @Failure      400           {object}  map[string]string                  "Invalid input"
// @Failure      500           {object}  map[string]string                  "Internal server error"
// @Router       /student/scenarios/phrase/listen [post]
func (h *PhraseStreamHandler) CreatePhraseStream(c *gin.Context) {
	var newPhraseStream models.CreatePhraseStreamRequest
	if err := c.ShouldBindJSON(&newPhraseStream); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	scenarioID, err := uuid.Parse(newPhraseStream.ScenarioID)
	phraseID, err := uuid.Parse(newPhraseStream.PhraseID)
	id, err := h.phraseStreamService.CreatePhraseStream(&domain.PhraseStream{
		ScenarioID: scenarioID,
		PhraseID:   phraseID,
		Status:     "initialized",
	}, &domain.AudioPhrase{
		PathToAudio: newPhraseStream.Path,
		PhraseID:    phraseID,
		Accent:      newPhraseStream.Accent,
		Noise:       newPhraseStream.Noise,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, id)
}

func (h *PhraseStreamHandler) AddAccent(c *gin.Context) {

}

// GetPhrases godoc
// @Summary      Get student phrases
// @Description  Returns a list of phrases associated with the given user
// @Tags         progress
// @Produce      json
// @Param        user_id  path      string  true  "User ID" Format(uuid)
// @Success      200      {array}   domain.Phrase
// @Failure      400      {object}  map[string]string  "Invalid user ID"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /student/{user_id}/get_phrases [get]
func (h *PhraseStreamHandler) GetPhrases(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	phrases, err := h.phraseStreamService.GetStudentPhrases(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, phrases)
}

// GetProgress godoc
// @Summary      Get phrase progress
// @Description  Returns progress data for the user's phrase practice
// @Tags         progress
// @Produce      json
// @Param        user_id  path      string  true  "User ID" Format(uuid)
// @Success      200      {array}   models.PhraseProgress
// @Failure      400      {object}  map[string]string  "Invalid user ID"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /student/{user_id}/phrase/get_progress [get]
func (h *PhraseStreamHandler) GetProgress(c *gin.Context) {
	userID := c.Param("user_id")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	phrases, err := h.phraseStreamService.GetStudentProgress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var result models.Progress
	for _, i := range phrases {
		result = append(result, &models.PhraseProgress{
			Phrase:             i[0],
			PhraseStreamStatus: i[1],
			ScenarioStatus:     i[2],
		})
	}
	c.JSON(http.StatusOK, phrases)
}
