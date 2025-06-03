package services

import (
	"diplom/internal/domain"
	"diplom/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) CreateUser(user *domain.User) (uuid.UUID, error) {
	return u.repo.Create(user)
}

func (u *UserService) GetUserByID(id uuid.UUID) (*domain.User, error) {
	return u.repo.GetByID(id)
}
func (u *UserService) IsAdmin(userID uuid.UUID) bool {
	user, err := u.GetUserByID(userID)
	if err != nil {
		return false
	}

	return user.Role
}

func (u *UserService) Login(login string, password string) (*domain.User, error) {
	user, err := u.repo.Login(login)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return user, err
}
