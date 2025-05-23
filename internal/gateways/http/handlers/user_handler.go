package handlers

import (
	"diplom/internal/domain"
	"diplom/internal/gateways/http/models"
	"diplom/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	user *services.UserService
}

func NewUserHandler(u *services.UserService) *UserHandler {
	return &UserHandler{u}
}

// RegisterUser godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.CreateUserRequest true "Данные для регистрации"
// @Success 201 {integer} integer "ID созданного пользователя"
// @Failure 400 {object} object "Неверный формат данных"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var newUser *models.CreateUserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := h.user.CreateUser(&domain.User{
		Name:     newUser.Name,
		Login:    newUser.Login,
		Password: newUser.Password,
		Role:     newUser.IsAdmin,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}
	c.JSON(http.StatusCreated, userID)
}

// LoginUser godoc
// @Summary Аутентификация пользователя
// @Description Выполняет вход пользователя в систему и возвращает информацию об успешной аутентификации
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginUserRequest true "Данные для входа (логин и пароль)"
//
//	@Success 200 {object} object "Успешная аутентификация" {
//	   "message"="Login successful",
//	   "user"="string"
//	}
//
//	@Failure 400 {object} object "Неверный формат запроса" {
//	   "error"="string"
//	}
//
//	@Failure 401 {object} object "Неверные учетные данные" {
//	   "message"="Invalid credentials"
//	}
//
// @Router /auth/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginUser models.LoginUserRequest
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.user.Login(loginUser.Login, loginUser.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user.Name})

}
