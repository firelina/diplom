package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AudioAnswerRepository struct {
	db *pgxpool.Pool
}

func NewAudioAnswerRepository(db *pgxpool.Pool) *AudioAnswerRepository {
	return &AudioAnswerRepository{db: db}
}

func (r *AudioAnswerRepository) Create(audioAnswer *domain.AudioAnswer) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.audio_answers (id, path_to_audio, record_time) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(context.Background(), query, id, audioAnswer.PathToAudio, audioAnswer.RecordTime)
	return id, err
}

func (r *AudioAnswerRepository) GetByID(id uuid.UUID) (*domain.AudioAnswer, error) {
	query := `SELECT id, path_to_audio, record_time FROM diplom.audio_answers WHERE id = $1`
	audioAnswer := &domain.AudioAnswer{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&audioAnswer.ID, &audioAnswer.PathToAudio, &audioAnswer.RecordTime)

	if err != nil {
		return nil, err
	}
	return audioAnswer, nil
}

//func (r *AudioAnswerRepository) Update(audioAnswer *domain.AudioAnswer) error {
//	query := `UPDATE audio_answers SET path_to_audio = $2, record_time = $3 WHERE id = $1`
//	_, err := r.db.Exec(context.Background(), query, audioAnswer.ID, audioAnswer.PathToAudio, audioAnswer.RecordTime)
//	return err
//}

func (r *AudioAnswerRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.audio_answers WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

//func (r *AudioAnswerRepository) GetAll() ([]domain.AudioAnswer, error) {
//	query := `SELECT id, path_to_audio, record_time FROM audio_answers`
//	rows, err := r.db.Query(context.Background(), query)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var audioAnswers []domain.AudioAnswer
//	for rows.Next() {
//		audioAnswer := domain.AudioAnswer{}
//		if err := rows.Scan(&audioAnswer.ID, &audioAnswer.PathToAudio, &audioAnswer.RecordTime); err != nil {
//			return nil, err
//		}
//		audioAnswers = append(audioAnswers, audioAnswer)
//	}
//	return audioAnswers, nil
//}
