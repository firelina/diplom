package services

import (
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockAudioPhraseRepository struct {
	mock.Mock
}

func (m *MockAudioPhraseRepository) Create(audioPhrase *domain.AudioPhrase) (uuid.UUID, error) {
	args := m.Called(audioPhrase)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAudioPhraseRepository) GetByID(id uuid.UUID) (*domain.AudioPhrase, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.AudioPhrase), args.Error(1)
}

func (m *MockAudioPhraseRepository) Update(audioPhrase *domain.AudioPhrase) error {
	args := m.Called(audioPhrase)
	return args.Error(0)
}

func (m *MockAudioPhraseRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAudioPhraseRepository) GetAll() ([]domain.AudioPhrase, error) {
	args := m.Called()
	return args.Get(0).([]domain.AudioPhrase), args.Error(1)
}

func TestAudioPhraseService_Create(t *testing.T) {
	mockRepo := new(MockAudioPhraseRepository)

	testID := uuid.New()
	phrase := &domain.AudioPhrase{
		PathToAudio: "/path/to/audio.mp3",
		PhraseID:    uuid.New(),
		Accent:      "British",
		Noise:       0.5,
	}

	mockRepo.On("Create", phrase).Return(testID, nil)

	id, err := mockRepo.Create(phrase)

	assert.NoError(t, err)
	assert.Equal(t, testID, id)
	mockRepo.AssertExpectations(t)
}
