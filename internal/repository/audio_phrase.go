package repository

import (
	"context"
	"diplom/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AudioPhraseRepository struct {
	db *pgxpool.Pool
}

func NewAudioPhraseRepository(db *pgxpool.Pool) *AudioPhraseRepository {
	return &AudioPhraseRepository{db: db}
}

func (r *AudioPhraseRepository) Create(audioPhrase *domain.AudioPhrase) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO diplom.audio_phrases (id, path_to_audio, phrase_id, accent, noise) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query, id, audioPhrase.PathToAudio, audioPhrase.PhraseID, audioPhrase.Accent, int(audioPhrase.Noise))
	return id, err
}

func (r *AudioPhraseRepository) GetByID(id uuid.UUID) (*domain.AudioPhrase, error) {
	query := `SELECT id, path_to_audio, phrase_id, accent, noise FROM diplom.audio_phrases WHERE id = $1`
	audioPhrase := &domain.AudioPhrase{}
	err := r.db.QueryRow(context.Background(), query, id).Scan(&audioPhrase.ID, &audioPhrase.PathToAudio, &audioPhrase.PhraseID, &audioPhrase.Accent, &audioPhrase.Noise)

	if err != nil {
		return nil, err
	}
	return audioPhrase, nil
}

func (r *AudioPhraseRepository) Update(audioPhrase *domain.AudioPhrase) error {
	query := `UPDATE diplom.audio_phrases SET path_to_audio = $2, phrase_id = $3, accent = $4, noise = $5 WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, audioPhrase.ID, audioPhrase.PathToAudio, audioPhrase.PhraseID, audioPhrase.Accent, audioPhrase.Noise)
	return err
}

func (r *AudioPhraseRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM diplom.audio_phrases WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *AudioPhraseRepository) GetAll() ([]domain.AudioPhrase, error) {
	query := `SELECT id, path_to_audio, phrase_id, accent, noise FROM diplom.audio_phrases`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audioPhrases []domain.AudioPhrase
	for rows.Next() {
		audioPhrase := domain.AudioPhrase{}
		if err := rows.Scan(&audioPhrase.ID, &audioPhrase.PathToAudio, &audioPhrase.PhraseID, &audioPhrase.Accent, &audioPhrase.Noise); err != nil {
			return nil, err
		}
		audioPhrases = append(audioPhrases, audioPhrase)
	}
	return audioPhrases, nil
}
