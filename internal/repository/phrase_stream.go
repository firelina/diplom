package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhraseStreamRepository struct {
	db *pgxpool.Pool
}

// NewPhraseStreamRepository initializes a new PhraseStreamRepository
func NewPhraseStreamRepository(db *pgxpool.Pool) *PhraseStreamRepository {
	return &PhraseStreamRepository{db: db}
}

// Create inserts a new PhraseStream record into the database
func (r *PhraseStreamRepository) Create(phraseStream *domain.PhraseStream) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.phrase_streams (id, audio_phrase_id, scenario_id, answer_id, phrase_id, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	_, err := r.db.Exec(context.Background(), query, id, phraseStream.AudioPhraseID, phraseStream.ScenarioID, "5d6629cb-ea31-42c6-8214-1732ab8619ea", phraseStream.PhraseID, phraseStream.Status)
	return id, err
}

// GetByID retrieves a PhraseStream record by its ID
func (r *PhraseStreamRepository) GetByID(id uuid.UUID) (*domain.PhraseStream, error) {
	query := `SELECT id, audio_phrase_id, scenario_id, answer_id, phrase_id, status FROM diplom.phrase_streams WHERE id = $1`
	phraseStream := &domain.PhraseStream{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&phraseStream.ID, &phraseStream.AudioPhraseID, &phraseStream.ScenarioID, &phraseStream.AnswerID, &phraseStream.PhraseID, &phraseStream.Status)

	if err != nil {
		return nil, err
	}
	return phraseStream, nil
}

func (r *PhraseStreamRepository) Update(id uuid.UUID, answerID uuid.UUID, status string) error {
	query := `UPDATE diplom.phrase_streams SET answer_id = $2, status = $3 WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id, answerID, status)
	return err
}

func (r *PhraseStreamRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.phrase_streams WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *PhraseStreamRepository) GetAll() ([]domain.PhraseStream, error) {
	query := `SELECT id, audio_phrase_id, scenario_id, answer_id, phrase_id, status FROM diplom.phrase_streams`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phraseStreams []domain.PhraseStream
	for rows.Next() {
		phraseStream := domain.PhraseStream{}
		if err := rows.Scan(&phraseStream.ID, &phraseStream.AudioPhraseID, &phraseStream.ScenarioID, &phraseStream.AnswerID, &phraseStream.PhraseID, &phraseStream.Status); err != nil {
			return nil, err
		}
		phraseStreams = append(phraseStreams, phraseStream)
	}
	return phraseStreams, nil
}

func (r *PhraseStreamRepository) GetStudentPhrases(userID uuid.UUID) ([]string, error) {
	query := `SELECT p.text FROM diplom.phrases p JOIN diplom.phrase_streams ps ON p.id = ps.phrase_id JOIN
    diplom.scenarios s on s.id = ps.scenario_id WHERE s.user_id =  $1`
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phrases []string
	for rows.Next() {
		var phrase string
		if err := rows.Scan(&phrase); err != nil {
			return nil, err
		}
		phrases = append(phrases, phrase)
	}
	return phrases, nil
}

func (r *PhraseStreamRepository) GetStudentProgress(userID uuid.UUID) ([][]string, error) {
	query := `SELECT p.text, ps.status, s.status FROM diplom.phrases p JOIN diplom.phrase_streams ps ON p.id = ps.phrase_id JOIN
    diplom.scenarios s on s.id = ps.scenario_id WHERE s.user_id =  $1`
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress [][]string
	for rows.Next() {
		var phrase string
		var phraseStreamStatus string
		var scenarioStatus string
		if err := rows.Scan(&phrase, &phraseStreamStatus, &scenarioStatus); err != nil {
			return nil, err
		}
		progress = append(progress, []string{phrase, phraseStreamStatus, scenarioStatus})
	}
	return progress, nil
}
