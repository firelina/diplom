package services

import (
	"diplom/internal/domain"
	"diplom/internal/repository"
	"github.com/google/uuid"
)

//type PhraseService interface {
//	CreatePhrase(phrase *domain.Phrase) error
//	GetPhraseByID(id uuid.UUID) (*domain.Phrase, error)
//	UpdatePhrase(phrase *domain.Phrase) error
//	DeletePhrase(id uuid.UUID) error
//	GetAllPhrases() ([]domain.Phrase, error)
//}

type PhraseService struct {
	repo repository.PhraseRepositoryInterface
}

func NewPhraseService(repo repository.PhraseRepositoryInterface) *PhraseService {
	return &PhraseService{repo: repo}
}

func (s *PhraseService) CreatePhrase(phrase *domain.Phrase) (uuid.UUID, error) {
	return s.repo.Create(phrase)
}

func (s *PhraseService) GetPhraseByID(id uuid.UUID) (*domain.Phrase, error) {
	return s.repo.GetByID(id)
}

func (s *PhraseService) UpdatePhrase(phrase *domain.Phrase) error {
	return s.repo.Update(phrase)
}

func (s *PhraseService) DeletePhrase(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *PhraseService) GetAllPhrases(text string) ([]domain.Phrase, error) {
	return s.repo.GetAll(text)
}
