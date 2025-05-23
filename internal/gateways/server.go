package gateways

import (
	"context"
	"diplom/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	host   string
	port   uint16
	router *gin.Engine
}

func (s *Server) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {

}

type Services struct {
	User         *services.UserService
	Phrase       *services.PhraseService
	PhraseType   *services.PhraseTypeService
	Answer       *services.StudentAnswerService
	Scenario     *services.ScenarioService
	PhraseStream *services.PhraseStreamService
}

func NewServer(services Services, options ...func(*Server)) *Server {
	r := gin.Default()

	setupRouter(r, services)

	s := &Server{router: r, host: "0.0.0.0", port: 8080}
	for _, o := range options {
		o(s)
	}

	return s
}

func (s *Server) Run(_ context.Context) error {
	return s.router.Run(fmt.Sprintf("%s:%d", s.host, s.port))
}
