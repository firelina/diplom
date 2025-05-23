package handlers

import (
	"diplom/internal/domain"
	"diplom/internal/gateways/http/models"
	"diplom/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type PhraseHandler struct {
	phraseService *services.PhraseService
	user          *services.UserService
}

func NewPhraseHandler(s *services.PhraseService, u *services.UserService) *PhraseHandler {
	return &PhraseHandler{phraseService: s, user: u}
}

// CreatePhrase godoc
// @Summary      Create a new phrase
// @Description  Adds a new phrase to the system
// @Tags         phrases
// @Accept       json
// @Produce      json
// @Param        phrase  body      models.CreatePhraseRequest  true  "New Phrase"
// @Success      201     {object}  string                       "Created ID"
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /admin/phrases [post]
func (h *PhraseHandler) CreatePhrase(c *gin.Context) {
	var newPhrase models.CreatePhraseRequest
	if err := c.ShouldBindJSON(&newPhrase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.phraseService.CreatePhrase(&domain.Phrase{
		Text:   newPhrase.Text,
		TypeID: newPhrase.TypeID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, id)
}

// GetPhrase godoc
// @Summary      Get a phrase by ID
// @Description  Returns a phrase by its UUID
// @Tags         phrases
// @Produce      json
// @Param        id   path      string           true  "Phrase ID" Format(uuid)
// @Success      200  {object}  domain.Phrase
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/phrases/{id} [get]
func (h *PhraseHandler) GetPhrase(c *gin.Context) {
	phraseID := c.Param("id")
	id, err := uuid.Parse(phraseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	phrase, err := h.phraseService.GetPhraseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Phrase not found"})
		return
	}
	c.JSON(http.StatusOK, phrase)
}

// UpdatePhrase godoc
// @Summary      Update a phrase by ID
// @Description  Updates an existing phrase
// @Tags         phrases
// @Accept       json
// @Produce      json
// @Param        id      path      string          true  "Phrase ID" Format(uuid)
// @Param        phrase  body      domain.Phrase   true  "Updated Phrase"
// @Success      200     {object}  domain.Phrase
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /admin/phrases/{id} [put]
func (h *PhraseHandler) UpdatePhrase(c *gin.Context) {
	phraseID := c.Param("id")
	id, err := uuid.Parse(phraseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	var updatedPhrase domain.Phrase
	if err := c.ShouldBindJSON(&updatedPhrase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedPhrase.ID = id
	if err := h.phraseService.UpdatePhrase(&updatedPhrase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedPhrase)
}

// DeletePhrase godoc
// @Summary      Delete a phrase by ID
// @Description  Deletes a phrase from the system
// @Tags         phrases
// @Produce      json
// @Param        id   path      string  true  "Phrase ID" Format(uuid)
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /admin/phrases/{id} [delete]
func (h *PhraseHandler) DeletePhrase(c *gin.Context) {
	phraseID := c.Param("id")
	id, err := uuid.Parse(phraseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := h.phraseService.DeletePhrase(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetAllPhrases godoc
// @Summary      Get all phrases
// @Description  Returns a list of all phrases
// @Tags         phrases
// @Produce      json
// @Success      200  {array}   domain.Phrase
// @Failure      500  {object}  map[string]string
// @Router       /admin/phrases [get]
func (h *PhraseHandler) GetAllPhrases(c *gin.Context) {
	phrases, err := h.phraseService.GetAllPhrases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, phrases)
}
