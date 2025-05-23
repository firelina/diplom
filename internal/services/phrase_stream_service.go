package services

import (
	"diplom/client"
	"diplom/internal/domain"
	"diplom/internal/repository"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/google/uuid"
	"github.com/jfreymuth/oggvorbis"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type PhraseStreamService struct {
	streams   *repository.PhraseStreamRepository
	audio     *repository.AudioPhraseRepository
	phrase    *repository.PhraseRepository
	speechKit *client.YandexSpeechClient
}

func NewPhraseStreamService(p *repository.PhraseStreamRepository, a *repository.AudioPhraseRepository, ph *repository.PhraseRepository) *PhraseStreamService {
	return &PhraseStreamService{
		streams:   p,
		audio:     a,
		phrase:    ph,
		speechKit: client.NewYandexSpeechClient(),
	}
}

func (s *PhraseStreamService) CreatePhraseStream(stream *domain.PhraseStream, audio *domain.AudioPhrase) (uuid.UUID, error) {
	pharse, err := s.phrase.GetByID(stream.PhraseID)
	if err != nil {
		return uuid.Nil, err
	}
	err = s.speechKit.SynthesizeSpeech(pharse.Text, audio.PathToAudio, audio.Accent)
	if err != nil {
		return uuid.Nil, err
	}
	err = addNoise(audio.PathToAudio, audio.Noise)
	if err != nil {
		return uuid.Nil, err
	}
	audioID, err := s.audio.Create(audio)
	if err != nil {
		return uuid.Nil, err
	}
	stream.AudioPhraseID = audioID
	return s.streams.Create(stream)
}

func (s *PhraseStreamService) UpdatePhraseStream(id uuid.UUID, answerID uuid.UUID, status string) error {
	return s.streams.Update(id, answerID, status)
}

func (s *PhraseStreamService) GetStudentPhrases(userID uuid.UUID) ([]string, error) {
	return s.streams.GetStudentPhrases(userID)
}

func (s *PhraseStreamService) GetStudentProgress(userID uuid.UUID) ([][]string, error) {
	return s.streams.GetStudentProgress(userID)
}

func addNoise(inputPath string, noiseLevel int) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer inFile.Close()

	decoder, err := oggvorbis.NewReader(inFile)
	if err != nil {
		return fmt.Errorf("ошибка при декодировании OGG: %w", err)
	}

	sampleRate := decoder.SampleRate()
	numChannels := decoder.Channels()

	var samples []float64
	buffer := make([]float32, 8192)

	for {
		n, err := decoder.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("ошибка при чтении сэмплов: %w", err)
		}
		for i := 0; i < n; i++ {
			samples = append(samples, float64(buffer[i]))
		}
	}

	rand.Seed(time.Now().UnixNano())
	for i := range samples {
		samples[i] += rand.NormFloat64() * float64(noiseLevel)
		if samples[i] > 1 {
			samples[i] = 1
		} else if samples[i] < -1 {
			samples[i] = -1
		}
	}

	intSamples := make([]int, len(samples))
	for i, s := range samples {
		intSamples[i] = int(s * (1 << 15))
	}

	buf := &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: numChannels,
			SampleRate:  sampleRate,
		},
		Data:           intSamples,
		SourceBitDepth: 16,
	}

	outputPath := inputPath[:len(inputPath)-len(filepath.Ext(inputPath))] + "_noisy.wav"
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("ошибка при создании WAV-файла: %w", err)
	}
	defer outFile.Close()

	encoder := wav.NewEncoder(outFile, sampleRate, 16, numChannels, 1)
	if err := encoder.Write(buf); err != nil {
		return fmt.Errorf("ошибка при записи WAV: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии WAV-файла: %w", err)
	}

	fmt.Println("WAV с шумом успешно сохранён в", outputPath)
	return nil
}
