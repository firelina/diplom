package services

import (
	"diplom/internal/domain"
	"diplom/internal/repository"
	"github.com/google/uuid"
)

type ScenarioService struct {
	repo *repository.ScenarioRepository
}

func NewScenarioService(repo *repository.ScenarioRepository) *ScenarioService {
	return &ScenarioService{repo: repo}
}

func (s *ScenarioService) CreateScenario(scenario *domain.Scenario) (uuid.UUID, error) {
	return s.repo.Create(scenario)
}

func (s *ScenarioService) GetScenario(id uuid.UUID) (*domain.Scenario, error) {
	return s.repo.GetByID(id)
}

func (s *ScenarioService) DeleteScenario(id uuid.UUID) error {
	return s.repo.Delete(id)
}
