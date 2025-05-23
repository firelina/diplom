package main

import (
	"bytes"
	"context"
	"diplom/internal/gateways"
	"diplom/internal/repository"
	"diplom/internal/services"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func generateFile() {
	apiKey := ""                                                           // Замените на ваш API-ключ
	text := "Requesting information on the weather conditions in the area" // Текст, который нужно преобразовать

	syn_url := "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"
	headers := map[string]string{
		"Authorization": "Api-Key " + apiKey,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	data := url.Values{}
	data.Set("text", text)
	data.Set("lang", "en-En")

	req, err := http.NewRequest("POST", syn_url, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	if resp.StatusCode == 200 {
		err = ioutil.WriteFile("output.ogg", body, 0644)
		if err != nil {
			fmt.Println("Ошибка при сохранении аудиофайла:", err)
			return
		}
		fmt.Println("Аудиофайл успешно создан: output.wav")
	} else {
		fmt.Println("Ошибка:", resp.Status, string(body))
	}
}

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
	//text := "Hello it's a text example on Go!"
	//
	//language := "en"
	//speed := "50"
	//volume := "100"
	//pitch := "50"
	//
	//outputFile := "output.ogg"
	//
	//cmd := exec.Command("espeak", "-v", language, "-s", speed, "-a", volume, "-p", pitch, "-w", outputFile, text)
	//
	//err := cmd.Run()
	//if err != nil {
	//	fmt.Println("Ошибка при преобразовании текста в речь:", err)
	//	return
	//}
	//
	//fmt.Println("Текст успешно преобразован в речь!")
	//generateFile()
	//// Путь к аудиофайлу
	//audioFilePath := "output.ogg"
	//
	//// Чтение аудиофайла
	//audioData, err := ioutil.ReadFile(audioFilePath)
	//if err != nil {
	//	fmt.Println("Ошибка при чтении файла:", err)
	//	return
	//}
	//
	//// URL для запроса
	//url_rec := "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize"
	//if err != nil {
	//	fmt.Println("Ошибка при создании запроса:", err)
	//	return
	//}
	//
	//req, err := http.NewRequest("POST", url_rec, bytes.NewBuffer(audioData))
	//
	//// Установка заголовков
	//req.Header.Set("Authorization", "Api-Key "+apiKey)
	//
	//// Отправка запроса
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("Ошибка при отправке запроса:", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//// Чтение ответа
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Ошибка при чтении ответа:", err)
	//	return
	//}
	//
	//// Проверка статуса ответа
	//if resp.StatusCode == 200 {
	//	fmt.Println("Распознанный текст:", string(body))
	//} else {
	//	fmt.Println("Ошибка:", resp.Status, string(body))
	//}
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
	phraseRepository := repository.NewPhraseRepository(pool)
	phraseTypeRepository := repository.NewPhraseTypeRepository(pool)
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
