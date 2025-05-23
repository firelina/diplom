package services

import (
	"diplom/internal/domain"
	"diplom/internal/repository"
	"github.com/google/uuid"
)

//type PhraseTypeService interface {
//	CreatePhraseType(phraseType *domain.PhraseType) error
//	GetPhraseTypeByID(id uuid.UUID) (*domain.PhraseType, error)
//	UpdatePhraseType(phraseType *domain.PhraseType) error
//	DeletePhraseType(id uuid.UUID) error
//	GetAllPhraseTypes() ([]domain.PhraseType, error)
//}

type PhraseTypeService struct {
	repo *repository.PhraseTypeRepository
}

func NewPhraseTypeService(repo *repository.PhraseTypeRepository) *PhraseTypeService {
	return &PhraseTypeService{repo: repo}
}

func (s *PhraseTypeService) CreatePhraseType(phraseType *domain.PhraseType) (uuid.UUID, error) {
	return s.repo.Create(phraseType)
}

func (s *PhraseTypeService) GetPhraseTypeByID(id uuid.UUID) (*domain.PhraseType, error) {
	return s.repo.GetByID(id)
}

func (s *PhraseTypeService) UpdatePhraseType(phraseType *domain.PhraseType) error {
	return s.repo.Update(phraseType)
}

func (s *PhraseTypeService) DeletePhraseType(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *PhraseTypeService) GetAllPhraseTypes() ([]domain.PhraseType, error) {
	return s.repo.GetAll()
}
