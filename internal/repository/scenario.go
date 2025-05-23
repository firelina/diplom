package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScenarioRepository struct {
	db *pgxpool.Pool
}

func NewScenarioRepository(db *pgxpool.Pool) *ScenarioRepository {
	return &ScenarioRepository{db: db}
}

func (r *ScenarioRepository) Create(scenario *domain.Scenario) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.scenarios (id, title, status, start_date, end_date, user_id) VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id`
	_, err := r.db.Exec(context.Background(), query, id, scenario.Title, scenario.Status, scenario.StartDate, scenario.EndDate, scenario.UserID)
	return id, err
}

func (r *ScenarioRepository) GetByID(id uuid.UUID) (*domain.Scenario, error) {
	query := `SELECT id, title, status, start_date, end_date, user_id FROM diplom.scenarios WHERE id = $1`
	scenario := &domain.Scenario{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&scenario.ID, &scenario.Title, &scenario.Status, &scenario.StartDate, &scenario.EndDate, &scenario.UserID)

	if err != nil {
		return nil, err
	}
	return scenario, nil
}

func (r *ScenarioRepository) Update(scenario *domain.Scenario) error {
	query := `UPDATE diplom.scenarios SET title = $2, status = $3, start_date = $4, end_date = $5, user_id = $6 WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, scenario.ID, scenario.Title, scenario.Status, scenario.StartDate, scenario.EndDate, scenario.UserID)
	return err
}

func (r *ScenarioRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.scenarios WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *ScenarioRepository) GetAll() ([]domain.Scenario, error) {
	query := `SELECT id, title, status, start_date, end_date, user_id FROM diplom.scenarios`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scenarios []domain.Scenario
	for rows.Next() {
		scenario := domain.Scenario{}
		if err := rows.Scan(&scenario.ID, &scenario.Title, &scenario.Status, &scenario.StartDate, &scenario.EndDate, &scenario.UserID); err != nil {
			return nil, err
		}
		scenarios = append(scenarios, scenario)
	}
	return scenarios, nil
}
