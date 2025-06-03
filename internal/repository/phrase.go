package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhraseRepositoryInterface interface {
	Create(phrase *domain.Phrase) (uuid.UUID, error)
	GetByID(id uuid.UUID) (*domain.Phrase, error)
	Update(phrase *domain.Phrase) error
	Delete(id uuid.UUID) error
	GetAll(textSearch string) ([]domain.Phrase, error)
}

type PhraseRepository struct {
	db *pgxpool.Pool
	pr PhraseTypeRepositoryInterface
}

func NewPhraseRepository(db *pgxpool.Pool, ph PhraseTypeRepositoryInterface) *PhraseRepository {
	return &PhraseRepository{db: db, pr: ph}
}

func (r *PhraseRepository) Create(phrase *domain.Phrase) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.phrases (id, text, type_id) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(context.Background(), query, id, phrase.Text, phrase.TypeID)
	return id, err
}

func (r *PhraseRepository) GetByID(id uuid.UUID) (*domain.Phrase, error) {
	query := `SELECT id, text, type_id FROM diplom.phrases WHERE id = $1`
	phrase := &domain.Phrase{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&phrase.ID, &phrase.Text, &phrase.TypeID)

	if err != nil {
		return nil, err
	}
	return phrase, nil
}

func (r *PhraseRepository) Update(phrase *domain.Phrase) error {
	query := `UPDATE diplom.phrases SET text = $2, type_id = $3 WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, phrase.ID, phrase.Text, phrase.TypeID)
	return err
}

func (r *PhraseRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.phrases WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *PhraseRepository) GetAll(textSearch string) ([]domain.Phrase, error) {
	query := `SELECT id, text, type_id FROM diplom.phrases WHERE text ILIKE '%' || $1 || '%'`
	rows, err := r.db.Query(context.Background(), query, textSearch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phrases []domain.Phrase
	for rows.Next() {
		phrase := domain.Phrase{}
		if err := rows.Scan(&phrase.ID, &phrase.Text, &phrase.TypeID); err != nil {
			return nil, err
		}
		phType, err := r.pr.GetByID(phrase.TypeID)
		if err != nil {
			return nil, err
		}
		phrase.PhraseType = phType.Title
		phrases = append(phrases, phrase)
	}
	return phrases, nil
}
