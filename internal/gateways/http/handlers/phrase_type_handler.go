package handlers

import (
	"diplom/internal/domain"
	"diplom/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type PhraseTypeHandler struct {
	phraseTypeService *services.PhraseTypeService
	user              *services.UserService
}

func NewPhraseTypeHandler(s *services.PhraseTypeService, u *services.UserService) *PhraseTypeHandler {
	return &PhraseTypeHandler{phraseTypeService: s, user: u}
}

// CreatePhraseType godoc
// @Summary      Create new phrase type
// @Description  Creates a new phrase type
// @Tags         phrase_types
// @Accept       json
// @Produce      json
// @Param        phrase_type  body      domain.PhraseType  true  "New Phrase Type"
// @Success      201          {object}  int                 "ID of created phrase type"
// @Failure      400          {object}  map[string]string   "Invalid input"
// @Failure      500          {object}  map[string]string   "Internal error"
// @Router       /admin/phrase_types [post]
func (h *PhraseTypeHandler) CreatePhraseType(c *gin.Context) {
	var newPhraseType domain.PhraseType
	if err := c.ShouldBindJSON(&newPhraseType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.phraseTypeService.CreatePhraseType(&newPhraseType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, id)
}

// GetAllPhraseTypes godoc
// @Summary      Get all phrase types
// @Description  Returns all phrase types
// @Tags         phrase_types
// @Produce      json
// @Success      200  {array}   domain.PhraseType
// @Failure      500  {object}  map[string]string  "Internal error"
// @Router       /admin/phrase_types [get]
func (h *PhraseTypeHandler) GetAllPhraseTypes(c *gin.Context) {
	phraseTypes, err := h.phraseTypeService.GetAllPhraseTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, phraseTypes)
}

func (h *PhraseTypeHandler) GetPhraseType(c *gin.Context) {
	phraseTypeID := c.Param("id")
	id, err := uuid.Parse(phraseTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	phraseType, err := h.phraseTypeService.GetPhraseTypeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phrase type not found"})
		return
	}
	c.JSON(http.StatusOK, phraseType)
}

func (h *PhraseTypeHandler) UpdatePhraseType(c *gin.Context) {
	phraseTypeID := c.Param("id")
	id, err := uuid.Parse(phraseTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	var updatedPhraseType domain.PhraseType
	if err := c.ShouldBindJSON(&updatedPhraseType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedPhraseType.ID = id
	if err := h.phraseTypeService.UpdatePhraseType(&updatedPhraseType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedPhraseType)
}

func (h *PhraseTypeHandler) DeletePhraseType(c *gin.Context) {
	phraseTypeID := c.Param("id")
	id, err := uuid.Parse(phraseTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := h.phraseTypeService.DeletePhraseType(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
