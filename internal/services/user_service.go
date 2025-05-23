package services

import (
	"diplom/internal/domain"
	"diplom/internal/repository"
	"github.com/google/uuid"
)

//type UserService interface {
//	CreateUser(user *domain.User) error
//	GetUserByID(id uuid.UUID) (*domain.User, error)
//	Login(login string, password string) (*domain.User, error)
//	IsAdmin(userID uuid.UUID) bool
//}

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
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
	return u.repo.Login(login, password)
}
