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

type StudentAnswerHandler struct {
	studentAnswerService *services.StudentAnswerService
	user                 *services.UserService
}

func NewStudentAnswerHandler(s *services.StudentAnswerService, u *services.UserService) *StudentAnswerHandler {
	return &StudentAnswerHandler{studentAnswerService: s, user: u}
}

// CreateAnswer godoc
// @Summary      Create a student answer
// @Description  Saves a student's audio answer to a phrase stream
// @Tags         scenarios
// @Accept       json
// @Produce      json
// @Param        answer  body      models.CreateAnswerRequest  true  "Student audio answer data"
// @Success      201     {object}  string                      "Created answer ID"
// @Failure      400     {object}  map[string]string           "Invalid request"
// @Failure      500     {object}  map[string]string           "Internal server error"
// @Router       /student/scenarios/answer [post]
func (h *StudentAnswerHandler) CreateAnswer(c *gin.Context) {
	var newAnswer models.CreateAnswerRequest

	if err := c.ShouldBindJSON(&newAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := uuid.Parse(newAnswer.UserID)
	phraseStreamID, err := uuid.Parse(newAnswer.PhraseStreamID)
	recordTime := time.Now()

	id, isCorrect, err := h.studentAnswerService.CreateAnswer(&domain.Answer{
		UserID: userID,
	}, &domain.AudioAnswer{
		PathToAudio: newAnswer.Path,
		RecordTime:  recordTime,
	}, phraseStreamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.CreateAnswerResponse{AnswerID: id, IsCorrect: isCorrect})
}

func (h *StudentAnswerHandler) GetAnswer(c *gin.Context) {
	answerID := c.Param("id")
	id, err := uuid.Parse(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	answer, audio, err := h.studentAnswerService.GetAnswer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"answer": answer, "audio": audio})
}

func (h *StudentAnswerHandler) UpdateAnswer(c *gin.Context) {
	answerID := c.Param("id")
	id, err := uuid.Parse(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	var updatedAnswer domain.Answer
	if err := c.ShouldBindJSON(&updatedAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedAnswer.ID = id
	if err := h.studentAnswerService.UpdateAnswer(&updatedAnswer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedAnswer)
}

// DeleteAnswer godoc
// @Summary      Delete a student answer
// @Description  Deletes a student answer by its ID
// @Tags         answers
// @Produce      json
// @Param        id   path      string  true  "Answer ID" Format(uuid)
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /admin/answers/{id} [delete]
func (h *StudentAnswerHandler) DeleteAnswer(c *gin.Context) {
	answerID := c.Param("id")
	id, err := uuid.Parse(answerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := h.studentAnswerService.DeleteAnswer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetAllAnswers godoc
// @Summary      Get all student answers
// @Description  Returns a list of all student answers
// @Tags         answers
// @Produce      json
// @Success      200  {array}   domain.Answer
// @Failure      500  {object}  map[string]string
// @Router       /admin/answers [get]
func (h *StudentAnswerHandler) GetAllAnswers(c *gin.Context) {
	answers, err := h.studentAnswerService.GetAllAnswers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, answers)
}
