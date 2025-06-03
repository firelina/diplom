package services

import (
	"diplom/internal/domain"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockPhraseStreamRepository struct {
	mock.Mock
}

func (m *MockPhraseStreamRepository) Create(phraseStream *domain.PhraseStream) (uuid.UUID, error) {
	args := m.Called(phraseStream)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockPhraseStreamRepository) GetByID(id uuid.UUID) (*domain.PhraseStream, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.PhraseStream), args.Error(1)
}

func (m *MockPhraseStreamRepository) Update(id uuid.UUID, answerID uuid.UUID, status string) error {
	args := m.Called(id, answerID, status)
	return args.Error(0)
}

func (m *MockPhraseStreamRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPhraseStreamRepository) GetAll() ([]domain.PhraseStream, error) {
	args := m.Called()
	return args.Get(0).([]domain.PhraseStream), args.Error(1)
}

func (m *MockPhraseStreamRepository) GetStudentPhrases(userID uuid.UUID) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPhraseStreamRepository) GetStudentProgress(userID uuid.UUID) ([][]string, error) {
	args := m.Called(userID)
	return args.Get(0).([][]string), args.Error(1)
}

func TestUpdatePhraseStream(t *testing.T) {
	mockRepo := new(MockPhraseStreamRepository)
	phraseMockRepo := new(MockPhraseRepository)
	audioPhraseMock := new(MockAudioPhraseRepository)
	service := NewPhraseStreamService(mockRepo, audioPhraseMock, phraseMockRepo)

	t.Run("success update phrase stream", func(t *testing.T) {
		id := uuid.New()
		answerID := uuid.New()
		status := "completed"

		mockRepo.On("Update", id, answerID, status).Return(nil)

		err := service.UpdatePhraseStream(id, answerID, status)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail update phrase stream", func(t *testing.T) {
		id := uuid.New()
		answerID := uuid.New()
		status := "completed"

		mockRepo.On("Update", id, answerID, status).Return(errors.New("update error"))

		err := service.UpdatePhraseStream(id, answerID, status)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetStudentPhrases(t *testing.T) {
	mockRepo := new(MockPhraseStreamRepository)
	phraseMockRepo := new(MockPhraseRepository)
	audioPhraseMock := new(MockAudioPhraseRepository)
	service := NewPhraseStreamService(mockRepo, audioPhraseMock, phraseMockRepo)

	t.Run("success get student phrases", func(t *testing.T) {
		userID := uuid.New()
		expectedPhrases := []string{"phrase1", "phrase2"}

		mockRepo.On("GetStudentPhrases", userID).Return(expectedPhrases, nil)

		phrases, err := service.GetStudentPhrases(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPhrases, phrases)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail get student phrases", func(t *testing.T) {
		userID := uuid.New()

		mockRepo.On("GetStudentPhrases", userID).Return([]string{}, errors.New("get error"))

		_, err := service.GetStudentPhrases(userID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetStudentProgress(t *testing.T) {
	mockRepo := new(MockPhraseStreamRepository)
	phraseMockRepo := new(MockPhraseRepository)
	audioPhraseMock := new(MockAudioPhraseRepository)
	service := NewPhraseStreamService(mockRepo, audioPhraseMock, phraseMockRepo)

	t.Run("success get student progress", func(t *testing.T) {
		userID := uuid.New()
		expectedProgress := [][]string{
			{"phrase1", "completed", "active"},
			{"phrase2", "pending", "inactive"},
		}

		mockRepo.On("GetStudentProgress", userID).Return(expectedProgress, nil)

		progress, err := service.GetStudentProgress(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedProgress, progress)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail get student progress", func(t *testing.T) {
		userID := uuid.New()

		mockRepo.On("GetStudentProgress", userID).Return([][]string{}, errors.New("get error"))

		_, err := service.GetStudentProgress(userID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
