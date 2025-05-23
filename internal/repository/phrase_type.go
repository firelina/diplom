package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhraseTypeRepository struct {
	db *pgxpool.Pool
}

func NewPhraseTypeRepository(db *pgxpool.Pool) *PhraseTypeRepository {
	return &PhraseTypeRepository{db: db}
}

func (r *PhraseTypeRepository) Create(phraseType *domain.PhraseType) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.phrase_types (id, title) VALUES ($1, $2)`
	_, err := r.db.Exec(context.Background(), query, id, phraseType.Title)
	return id, err
}

func (r *PhraseTypeRepository) GetByID(id uuid.UUID) (*domain.PhraseType, error) {
	query := `SELECT id, title FROM diplom.phrase_types WHERE id = $1`
	phraseType := &domain.PhraseType{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&phraseType.ID, &phraseType.Title)

	if err != nil {
		return nil, err
	}
	return phraseType, nil
}

func (r *PhraseTypeRepository) Update(phraseType *domain.PhraseType) error {
	query := `UPDATE diplom.phrase_types SET title = $2 WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, phraseType.ID, phraseType.Title)
	return err
}

func (r *PhraseTypeRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.phrase_types WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *PhraseTypeRepository) GetAll() ([]domain.PhraseType, error) {
	query := `SELECT id, title FROM diplom.phrase_types`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phraseTypes []domain.PhraseType
	for rows.Next() {
		phraseType := domain.PhraseType{}
		if err := rows.Scan(&phraseType.ID, &phraseType.Title); err != nil {
			return nil, err
		}
		phraseTypes = append(phraseTypes, phraseType)
	}
	return phraseTypes, nil
}
