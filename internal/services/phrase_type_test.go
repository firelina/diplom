package services

import (
	"diplom/internal/domain"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockPhraseTypeRepository struct {
	mock.Mock
}

func (m *MockPhraseTypeRepository) Create(phraseType *domain.PhraseType) (uuid.UUID, error) {
	args := m.Called(phraseType)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockPhraseTypeRepository) GetByID(id uuid.UUID) (*domain.PhraseType, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.PhraseType), args.Error(1)
}

func (m *MockPhraseTypeRepository) Update(phraseType *domain.PhraseType) error {
	args := m.Called(phraseType)
	return args.Error(0)
}

func (m *MockPhraseTypeRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPhraseTypeRepository) GetAll() ([]domain.PhraseType, error) {
	args := m.Called()
	return args.Get(0).([]domain.PhraseType), args.Error(1)
}

func TestPhraseTypeService_CreatePhraseType(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		testID := uuid.New()
		phraseType := &domain.PhraseType{Title: "Test Type"}

		mockRepo.On("Create", phraseType).Return(testID, nil)

		id, err := srv.CreatePhraseType(phraseType)

		assert.NoError(t, err)
		assert.Equal(t, testID, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("creation error", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		expectedErr := errors.New("repository error")
		phraseType := &domain.PhraseType{Title: "Test Type"}

		mockRepo.On("Create", phraseType).Return(uuid.Nil, expectedErr)

		id, err := srv.CreatePhraseType(phraseType)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, uuid.Nil, id)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseTypeService_GetPhraseTypeByID(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		testID := uuid.New()
		expectedType := &domain.PhraseType{
			ID:    testID,
			Title: "Test Type",
		}

		mockRepo.On("GetByID", testID).Return(expectedType, nil)

		result, err := srv.GetPhraseTypeByID(testID)

		assert.NoError(t, err)
		assert.Equal(t, expectedType, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseTypeService_UpdatePhraseType(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		phraseType := &domain.PhraseType{
			ID:    uuid.New(),
			Title: "Updated Type",
		}

		mockRepo.On("Update", phraseType).Return(nil)

		err := srv.UpdatePhraseType(phraseType)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		expectedErr := errors.New("update error")
		phraseType := &domain.PhraseType{
			ID:    uuid.New(),
			Title: "Updated Type",
		}

		mockRepo.On("Update", phraseType).Return(expectedErr)

		err := srv.UpdatePhraseType(phraseType)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseTypeService_DeletePhraseType(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		testID := uuid.New()

		mockRepo.On("Delete", testID).Return(nil)

		err := srv.DeletePhraseType(testID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		testID := uuid.New()
		expectedErr := errors.New("delete error")

		mockRepo.On("Delete", testID).Return(expectedErr)

		err := srv.DeletePhraseType(testID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPhraseTypeService_GetAllPhraseTypes(t *testing.T) {
	t.Run("successful get all", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		expectedTypes := []domain.PhraseType{
			{ID: uuid.New(), Title: "Type 1"},
			{ID: uuid.New(), Title: "Type 2"},
		}

		mockRepo.On("GetAll").Return(expectedTypes, nil)

		result, err := srv.GetAllPhraseTypes()

		assert.NoError(t, err)
		assert.Equal(t, expectedTypes, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get all error", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		expectedErr := errors.New("repository error")

		mockRepo.On("GetAll").Return([]domain.PhraseType{}, expectedErr)

		result, err := srv.GetAllPhraseTypes()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(MockPhraseTypeRepository)
		srv := NewPhraseTypeService(mockRepo)

		mockRepo.On("GetAll").Return([]domain.PhraseType{}, nil)

		result, err := srv.GetAllPhraseTypes()

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}
