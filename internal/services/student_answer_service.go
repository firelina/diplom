package services

import (
	"diplom/client"
	"diplom/internal/domain"
	"diplom/internal/repository"
	"github.com/google/uuid"
	"math"
	"strings"
)

//type StudentAnswerService interface {
//	CreateAnswer(answer *domain.Answer, audio *domain.AudioAnswer) error
//	CreateAudioAnswer(answer *domain.AudioAnswer) (uuid.UUID, error)
//	GetAllAnswers() ([]domain.Answer, error)
//	GetAnswerByID(id uuid.UUID) (*domain.Answer, error)
//	GetAudioAnswerByID(id uuid.UUID) (*domain.AudioAnswer, error)
//	GetAnswer(id uuid.UUID) (*domain.Answer, *domain.AudioAnswer, error)
//	UpdateAnswer(answer *domain.Answer) error
//	DeleteAnswer(id uuid.UUID) error
//	DeleteAudioAnswer(id uuid.UUID) error
//}

type StudentAnswerService struct {
	answerRepository      *repository.AnswerRepository
	audioAnswerRepository *repository.AudioAnswerRepository

	phraseStream *repository.PhraseStreamRepository
	phrase       *repository.PhraseRepository
	speechKit    *client.YandexSpeechClient
}

func NewStudentAnswerService(answer *repository.AnswerRepository, audio *repository.AudioAnswerRepository,
	phs *repository.PhraseStreamRepository, ph *repository.PhraseRepository) *StudentAnswerService {
	return &StudentAnswerService{answerRepository: answer, audioAnswerRepository: audio,
		phraseStream: phs, phrase: ph, speechKit: client.NewYandexSpeechClient()}
}

func (s *StudentAnswerService) CreateAnswer(answer *domain.Answer, audio *domain.AudioAnswer, phraseStreamID uuid.UUID) (uuid.UUID, bool, string, error) {
	phraseStream, err := s.phraseStream.GetByID(phraseStreamID)
	if err != nil {
		return uuid.Nil, false, "", err
	}
	phrase, err := s.phrase.GetByID(phraseStream.PhraseID)
	if err != nil {
		return uuid.Nil, false, "", err
	}

	text, err := s.speechKit.RecognizeSpeech(audio.PathToAudio)
	similirity := CosineSimilarity(phrase.Text, text)
	if err != nil {
		return uuid.Nil, false, "", err
	}

	var status string
	var isCorrect bool
	if similirity > 0.1 {
		isCorrect = true
		status = "success"
	} else {
		isCorrect = false
		status = "fail"
	}

	audioID, err := s.audioAnswerRepository.Create(audio)
	if err != nil {
		return uuid.Nil, false, "", err
	}
	answer.AudioAnswerID = audioID
	answer.Text = text
	answer.IsCorrect = isCorrect
	answerID, err := s.answerRepository.Create(answer)
	if err != nil {
		return uuid.Nil, false, "", err
	}

	err = s.phraseStream.Update(phraseStreamID, answerID, status)
	if err != nil {
		return uuid.Nil, false, "", err
	}
	return answerID, isCorrect, text, err
}

func (s *StudentAnswerService) CreateAudioAnswer(answer *domain.AudioAnswer) (uuid.UUID, error) {
	//return s.audioAnswerRepository.Create(answer), nil
	return uuid.Nil, nil
}

func (s *StudentAnswerService) GetAllAnswers() ([]domain.Answer, error) {
	return s.answerRepository.GetAll()
}

func (s *StudentAnswerService) GetAnswerByID(id uuid.UUID) (*domain.Answer, error) {
	return s.answerRepository.GetByID(id)
}

func (s *StudentAnswerService) GetAudioAnswerByID(id uuid.UUID) (*domain.AudioAnswer, error) {
	//return s.audioAnswerRepository.GetByID(id), nil
	return nil, nil
}

func (s *StudentAnswerService) GetAnswer(id uuid.UUID) (*domain.Answer, *domain.AudioAnswer, error) {
	answer, err := s.GetAnswerByID(id)
	if err != nil {
		return nil, nil, err
	}
	audio, err := s.GetAudioAnswerByID(answer.AudioAnswerID)
	if err != nil {
		return nil, nil, err
	}
	return answer, audio, nil
}

func (s *StudentAnswerService) UpdateAnswer(answer *domain.Answer) error {
	return s.answerRepository.Update(answer)
}

func (s *StudentAnswerService) DeleteAnswer(id uuid.UUID) error {
	answer, err := s.GetAnswerByID(id)
	if err != nil {
		return err
	}
	err = s.DeleteAudioAnswer(answer.AudioAnswerID)
	if err != nil {
		return err
	}
	return s.answerRepository.Delete(id)
}

func (s *StudentAnswerService) DeleteAudioAnswer(id uuid.UUID) error {
	//return s.audioAnswerRepository.Delete(id)
	return nil
}

func tokenize(s string) []string {
	words := strings.Fields(strings.ToLower(s))
	return words
}

func termFrequency(tokens []string) map[string]float64 {
	tf := make(map[string]float64)
	for _, token := range tokens {
		tf[token]++
	}
	return tf
}

func CosineSimilarity(a, b string) float64 {
	tokensA := tokenize(a)
	tokensB := tokenize(b)

	tfA := termFrequency(tokensA)
	tfB := termFrequency(tokensB)

	vocab := make(map[string]bool)
	for word := range tfA {
		vocab[word] = true
	}
	for word := range tfB {
		vocab[word] = true
	}

	var dotProduct, normA, normB float64
	for word := range vocab {
		valA := tfA[word]
		valB := tfB[word]
		dotProduct += valA * valB
		normA += valA * valA
		normB += valB * valB
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
