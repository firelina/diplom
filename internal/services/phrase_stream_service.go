package services

import (
	"diplom/client"
	"diplom/internal/domain"
	"diplom/internal/repository"
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
	"io"
	"math/rand"
	"os"
	"time"

	//"fmt"
	"github.com/google/uuid"
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

func addNoise(inputPath string, noiseLevel float64) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия MP3: %w", err)
	}
	defer f.Close()

	decoder, err := mp3.NewDecoder(f)
	if err != nil {
		return fmt.Errorf("ошибка создания декодера: %w", err)
	}

	pcmBytes, err := io.ReadAll(decoder)
	if err != nil {
		return fmt.Errorf("ошибка чтения PCM данных: %w", err)
	}

	pcm := make([]int16, len(pcmBytes)/2)
	for i := 0; i < len(pcm); i++ {
		pcm[i] = int16(pcmBytes[2*i]) | int16(pcmBytes[2*i+1])<<8
	}

	rand.Seed(time.Now().UnixNano())
	maxInt16 := 32767
	for i := range pcm {
		noise := int(float64(maxInt16) * noiseLevel * rand.NormFloat64())
		val := int(pcm[i]) + noise
		pcm[i] = int16(clamp(val, -maxInt16, maxInt16))
	}

	outputWav := "output1.wav"
	outFile, err := os.Create(outputWav)
	if err != nil {
		return fmt.Errorf("ошибка создания WAV файла: %w", err)
	}
	defer outFile.Close()

	enc := wav.NewEncoder(outFile, 44100, 16, 2, 1) // Частота и каналы нужно брать из mp3, упрощено 44.1kHz, стерео
	defer enc.Close()

	buf := &audio.IntBuffer{
		Data:           make([]int, len(pcm)),
		Format:         &audio.Format{NumChannels: 2, SampleRate: 44100},
		SourceBitDepth: 16,
	}
	for i, v := range pcm {
		buf.Data[i] = int(v)
	}

	if err := enc.Write(buf); err != nil {
		return fmt.Errorf("ошибка записи WAV: %w", err)
	}

	//tempWav := "temp_output.wav"
	//outFile, err := os.Create(tempWav)
	//if err != nil {
	//	return fmt.Errorf("ошибка создания WAV файла: %w", err)
	//}
	//defer outFile.Close()
	//
	//enc := wav.NewEncoder(outFile, 44100, 16, 2, 1) // можно уточнить параметры у decoder
	//defer enc.Close()
	//
	//intBuf := &audio.IntBuffer{
	//	Data:           make([]int, len(pcm)),
	//	Format:         &audio.Format{SampleRate: 44100, NumChannels: 2},
	//	SourceBitDepth: 16,
	//}
	//for i, v := range pcm {
	//	intBuf.Data[i] = int(v)
	//}
	//if err := enc.Write(intBuf); err != nil {
	//	return fmt.Errorf("ошибка записи WAV: %w", err)
	//}
	//
	//// 4. Конвертируем в MP3
	//cmd := exec.Command("lame", tempWav, "output.mp3")
	//if out, err := cmd.CombinedOutput(); err != nil {
	//	return fmt.Errorf("ошибка при вызове lame: %w (%s)", err, string(out))
	//}
	//
	//// Удаляем временный WAV
	//_ = os.Remove(tempWav)
	return nil
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
