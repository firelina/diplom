package handlers

import (
	"diplom/internal/domain"
	"diplom/internal/gateways/http/models"
	"diplom/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ScenarioHandler struct {
	scenarioService *services.ScenarioService
}

func NewScenarioHandler(s *services.ScenarioService) *ScenarioHandler {
	return &ScenarioHandler{scenarioService: s}
}

// CreateScenario godoc
// @Summary      Create a new scenario
// @Description  Creates a new scenario for a student
// @Tags         scenarios
// @Accept       json
// @Produce      json
// @Param        scenario  body      models.CreateScenarioRequest  true  "Scenario data"
// @Success      201       {object}  string                        "Created scenario ID"
// @Failure      400       {object}  map[string]string             "Invalid input"
// @Failure      500       {object}  map[string]string             "Internal server error"
// @Router       /student/scenarios/create [post]
func (h *ScenarioHandler) CreateScenario(c *gin.Context) {
	var newScenario models.CreateScenarioRequest
	if err := c.ShouldBindJSON(&newScenario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := uuid.Parse(newScenario.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	startDate := time.Now()
	id, err := h.scenarioService.CreateScenario(&domain.Scenario{
		Title:     newScenario.Title,
		UserID:    userID,
		StartDate: &startDate,
		Status:    "in_progress",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, id)
}

func (h *ScenarioHandler) GetScenario(c *gin.Context) {
	scenarioID := c.Param("id")
	id, err := uuid.Parse(scenarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	scenario, err := h.scenarioService.GetScenario(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scenario not found"})
		return
	}
	c.JSON(http.StatusOK, scenario)
}

func (h *ScenarioHandler) DeleteScenario(c *gin.Context) {
	scenarioID := c.Param("id")
	id, err := uuid.Parse(scenarioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := h.scenarioService.DeleteScenario(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
