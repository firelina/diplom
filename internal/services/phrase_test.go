package services

import (
	"diplom/internal/domain"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockPhraseRepository struct {
	mock.Mock
	pr *MockPhraseTypeRepository
}

func (m *MockPhraseRepository) Create(phrase *domain.Phrase) (uuid.UUID, error) {
	args := m.Called(phrase)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockPhraseRepository) GetByID(id uuid.UUID) (*domain.Phrase, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Phrase), args.Error(1)
}

func (m *MockPhraseRepository) Update(phrase *domain.Phrase) error {
	args := m.Called(phrase)
	return args.Error(0)
}

func (m *MockPhraseRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPhraseRepository) GetAll(textSearch string) ([]domain.Phrase, error) {
	args := m.Called(textSearch)
	return args.Get(0).([]domain.Phrase), args.Error(1)
}

func TestPhraseRepository_Create(t *testing.T) {
	mockTypeRepo := new(MockPhraseTypeRepository)

	mockRepo := new(MockPhraseRepository)
	mockRepo.pr = mockTypeRepo
	service := NewPhraseService(mockRepo)

	t.Run("success create phrase", func(t *testing.T) {
		phrase := &domain.Phrase{
			Text:   "Test phrase",
			TypeID: uuid.New(),
		}
		expectedID := uuid.New()

		mockRepo.On("Create", phrase).Return(expectedID, nil)

		id, err := service.CreatePhrase(phrase)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail create phrase", func(t *testing.T) {
		phrase := &domain.Phrase{}
		mockRepo.On("Create", phrase).Return(uuid.Nil, errors.New("create error"))

		_, err := service.CreatePhrase(phrase)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseRepository_GetByID(t *testing.T) {
	mockTypeRepo := new(MockPhraseTypeRepository)

	mockRepo := new(MockPhraseRepository)
	mockRepo.pr = mockTypeRepo
	service := NewPhraseService(mockRepo)

	t.Run("success get phrase by id", func(t *testing.T) {
		id := uuid.New()
		expectedPhrase := &domain.Phrase{
			ID:     id,
			Text:   "Test phrase",
			TypeID: uuid.New(),
		}

		mockRepo.On("GetByID", id).Return(expectedPhrase, nil)

		phrase, err := service.GetPhraseByID(id)

		assert.NoError(t, err)
		assert.Equal(t, expectedPhrase, phrase)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail get phrase by id", func(t *testing.T) {
		id := uuid.New()
		mockRepo.On("GetByID", id).Return(&domain.Phrase{}, errors.New("not found"))

		_, err := service.GetPhraseByID(id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseRepository_Update(t *testing.T) {
	mockTypeRepo := new(MockPhraseTypeRepository)

	mockRepo := new(MockPhraseRepository)
	mockRepo.pr = mockTypeRepo
	service := NewPhraseService(mockRepo)

	t.Run("success update phrase", func(t *testing.T) {
		phrase := &domain.Phrase{
			ID:     uuid.New(),
			Text:   "Updated phrase",
			TypeID: uuid.New(),
		}

		mockRepo.On("Update", phrase).Return(nil)

		err := service.UpdatePhrase(phrase)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail update phrase", func(t *testing.T) {
		phrase := &domain.Phrase{}
		mockRepo.On("Update", phrase).Return(errors.New("update error"))

		err := service.UpdatePhrase(phrase)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseRepository_Delete(t *testing.T) {
	mockTypeRepo := new(MockPhraseTypeRepository)

	mockRepo := new(MockPhraseRepository)
	mockRepo.pr = mockTypeRepo
	service := NewPhraseService(mockRepo)

	t.Run("success delete phrase", func(t *testing.T) {
		id := uuid.New()

		mockRepo.On("Delete", id).Return(nil)

		err := service.DeletePhrase(id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail delete phrase", func(t *testing.T) {
		id := uuid.New()
		mockRepo.On("Delete", id).Return(errors.New("delete error"))

		err := service.DeletePhrase(id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseRepository_GetAll(t *testing.T) {
	mockTypeRepo := new(MockPhraseTypeRepository)

	mockRepo := new(MockPhraseRepository)
	mockRepo.pr = mockTypeRepo
	service := NewPhraseService(mockRepo)

	t.Run("success get all phrases", func(t *testing.T) {
		searchText := "test"
		expectedPhrases := []domain.Phrase{
			{ID: uuid.New(), Text: "Test phrase 1", TypeID: uuid.New()},
			{ID: uuid.New(), Text: "Test phrase 2", TypeID: uuid.New()},
		}

		mockRepo.On("GetAll", searchText).Return(expectedPhrases, nil)

		phrases, err := service.GetAllPhrases(searchText)

		assert.NoError(t, err)
		assert.Equal(t, expectedPhrases, phrases)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail get all phrases", func(t *testing.T) {
		searchText := ""
		mockRepo.On("GetAll", searchText).Return([]domain.Phrase{}, errors.New("get error"))

		_, err := service.GetAllPhrases(searchText)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
