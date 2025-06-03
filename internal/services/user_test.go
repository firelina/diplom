package services

import (
	"diplom/internal/domain"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) (uuid.UUID, error) {
	args := m.Called(user)
	id, _ := uuid.Parse("5d6629cb-ea31-42c6-8214-1732ab8619ea")
	return id, args.Error(1)
}

func (m *MockUserRepository) Login(login string) (*domain.User, error) {
	args := m.Called(login)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	return nil
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	return nil
}

func TestRegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	t.Run("success register user", func(t *testing.T) {
		newUser := &domain.User{Name: "John Doe", Login: "john@example.com", Password: "password123", Role: false}
		id, _ := uuid.Parse("5d6629cb-ea31-42c6-8214-1732ab8619ea")

		mockRepo.On("Create", newUser).Return(id, nil)

		userID, err := userService.CreateUser(newUser)

		assert.NoError(t, err)
		assert.Equal(t, id, userID)

		mockRepo.AssertExpectations(t)
	})

}

func TestLoginUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	t.Run("success login user", func(t *testing.T) {
		login := "john@example.com"
		password := "1234"
		expectedUser := &domain.User{ID: uuid.New(), Name: "John Doe", Login: login, Password: "$2a$10$q6O9J2b4BtFY224tmjC6.eAF6Keqz39/5Uu9aKtyKOpnNFLeeuCoC", Role: false}

		mockRepo.On("Login", login).Return(expectedUser, nil)

		user, err := userService.Login(login, password)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Login, user.Login)

		mockRepo.AssertExpectations(t)
	})

	t.Run("fail login user not found error", func(t *testing.T) {
		login := "john@example.com"
		password := "12345"
		expectedUser := &domain.User{ID: uuid.New(), Name: "John Doe", Login: login, Password: "$2a$10$q6O9J2b4BtFY224tmjC6.eAF6Keqz39/5Uu9aKtyKOpnNFLeeuCoC", Role: false}

		mockRepo.On("Login", login).Return(expectedUser, nil)

		_, err := userService.Login(login, password)

		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserIsAdmin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)
	t.Run("user is admin", func(t *testing.T) {
		newUser := &domain.User{Name: "John Doe", Login: "john@example.com", Password: "password123", Role: true}
		id, _ := uuid.Parse("5d6629cb-ea31-42c6-8214-1732ab8619ea")

		mockRepo.On("GetByID", id).Return(newUser, nil)

		role := userService.IsAdmin(id)

		assert.True(t, role)

		mockRepo.AssertExpectations(t)
	})

}
