package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnswerRepository struct {
	db *pgxpool.Pool
}

func NewAnswerRepository(db *pgxpool.Pool) *AnswerRepository {
	return &AnswerRepository{db: db}
}

func (r *AnswerRepository) Create(answer *domain.Answer) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.answers (id, user_id, audio_answer_id, text, is_correct) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	_, err := r.db.Exec(context.Background(), query, id, answer.UserID, answer.AudioAnswerID, answer.Text, answer.IsCorrect)
	return id, err
}

func (r *AnswerRepository) GetByID(id uuid.UUID) (*domain.Answer, error) {
	query := `SELECT id, user_id, audio_answer_id, text, is_correct FROM diplom.answers WHERE id = $1`
	answer := &domain.Answer{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&answer.ID, &answer.UserID, &answer.AudioAnswerID, &answer.Text, &answer.IsCorrect)

	if err != nil {
		return nil, err
	}
	return answer, nil
}

func (r *AnswerRepository) Update(answer *domain.Answer) error {
	query := `UPDATE diplom.answers SET user_id = $2, audio_answer_id = $3, text = $4, is_correct = $5 WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, answer.ID, answer.UserID, answer.AudioAnswerID, answer.Text, answer.IsCorrect)
	return err
}

func (r *AnswerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.answers WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *AnswerRepository) GetAll() ([]domain.Answer, error) {
	query := `SELECT id, user_id, audio_answer_id, text, is_correct FROM diplom.answers`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []domain.Answer
	for rows.Next() {
		answer := domain.Answer{}
		if err := rows.Scan(&answer.ID, &answer.UserID, &answer.AudioAnswerID, &answer.Text, &answer.IsCorrect); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}
	return answers, nil
}
