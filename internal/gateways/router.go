package gateways

import (
	"diplom/internal/gateways/http/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(r *gin.Engine, services Services) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.StaticFile("/swag/swagger.json", "./docs/swagger.json")
	url := ginSwagger.URL("http://localhost:8080/swag/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	userHandler := handlers.NewUserHandler(services.User)
	phraseTypeHandler := handlers.NewPhraseTypeHandler(services.PhraseType, services.User)
	phraseHandler := handlers.NewPhraseHandler(services.Phrase, services.User)
	answerHandler := handlers.NewStudentAnswerHandler(services.Answer, services.User)
	scenarioHandler := handlers.NewScenarioHandler(services.Scenario)
	phraseStreamHandler := handlers.NewPhraseStreamHandler(services.PhraseStream)

	r.POST("/api/v1/users/register", func(c *gin.Context) {
		userHandler.RegisterUser(c)
	})
	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		userHandler.LoginUser(c)
	})

	r.POST("/api/v1/admin/phrases", func(c *gin.Context) {
		phraseHandler.CreatePhrase(c)
	})
	r.GET("/api/v1/admin/phrases/:id", func(c *gin.Context) {
		phraseHandler.GetPhrase(c)
	})
	r.PUT("/api/v1/admin/phrases/:id", func(c *gin.Context) {
		phraseHandler.UpdatePhrase(c)
	})
	r.DELETE("/api/v1/admin/phrases/:id", func(c *gin.Context) {
		phraseHandler.DeletePhrase(c)
	})
	r.GET("/api/v1/admin/phrases", func(c *gin.Context) {
		phraseHandler.GetAllPhrases(c)
	})

	r.POST("/api/v1/admin/phrase_types", func(c *gin.Context) {
		phraseTypeHandler.CreatePhraseType(c)
	})
	r.GET("/api/v1/admin/phrase_types", func(c *gin.Context) {
		phraseTypeHandler.GetAllPhraseTypes(c)
	})

	r.GET("/api/v1/admin/answers", func(c *gin.Context) {
		answerHandler.GetAllAnswers(c)
	})
	r.DELETE("/api/v1/admin/answers/:id", func(c *gin.Context) {
		answerHandler.DeleteAnswer(c)
	})

	r.POST("/api/v1/student/scenarios/create", func(c *gin.Context) {
		scenarioHandler.CreateScenario(c)
	})
	r.POST("/api/v1/student/scenarios/answer", func(c *gin.Context) {
		answerHandler.CreateAnswer(c)
	})
	r.POST("/api/v1/student/scenarios/phrase/listen", func(c *gin.Context) {
		phraseStreamHandler.CreatePhraseStream(c)
	})
	r.GET("/api/v1/student/:user_id/get_phrases", func(c *gin.Context) {
		phraseStreamHandler.GetPhrases(c)
	})
	r.GET("/api/v1/student/:user_id/phrase/get_progress", func(c *gin.Context) {
		phraseStreamHandler.GetProgress(c)
	})
}
