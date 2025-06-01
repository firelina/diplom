package main

import (
	"context"
	"diplom/internal/gateways"
	"diplom/internal/repository"
	"diplom/internal/services"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title My Awesome API
// @version 1.0
// @description This is a sample server for a pet store.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	server := http.Server{
		ReadHeaderTimeout: 10 * time.Second,
	}
	config, err := pgxpool.ParseConfig("postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("can't parse pgxpool config")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("can't create new pool")
	}
	defer pool.Close()
	userRepository := repository.NewUserRepository(pool)
	phraseTypeRepository := repository.NewPhraseTypeRepository(pool)
	phraseRepository := repository.NewPhraseRepository(pool, phraseTypeRepository)
	answerRepository := repository.NewAnswerRepository(pool)
	phraseStreamRepository := repository.NewPhraseStreamRepository(pool)
	audioAnswerRepository := repository.NewAudioAnswerRepository(pool)
	audioPhraseRepository := repository.NewAudioPhraseRepository(pool)
	scenarioRepository := repository.NewScenarioRepository(pool)

	useCases := gateways.Services{
		User:         services.NewUserService(userRepository),
		Phrase:       services.NewPhraseService(phraseRepository),
		PhraseType:   services.NewPhraseTypeService(phraseTypeRepository),
		Answer:       services.NewStudentAnswerService(answerRepository, audioAnswerRepository, phraseStreamRepository, phraseRepository),
		Scenario:     services.NewScenarioService(scenarioRepository),
		PhraseStream: services.NewPhraseStreamService(phraseStreamRepository, audioPhraseRepository, phraseRepository),
	}
	r := gateways.NewServer(useCases)
	server.Handler = r

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	eg, _ := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		s := <-sigQuit
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Println(err.Error())
		}
		return fmt.Errorf("captured signal: %v", s)
	})

	go func() {
		if err := r.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error during server shutdown: %v", err)
		}
	}()
	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the server: %v", err) // gracefully shutting down the server: captured signal: interrupt
	}
}
