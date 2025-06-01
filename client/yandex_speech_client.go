package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type YandexSpeechClient struct {
	apiKey string
	client *http.Client
}

type Answer struct {
	Result string `json:"result"`
}

func NewYandexSpeechClient() *YandexSpeechClient {
	return &YandexSpeechClient{
		apiKey: "",
		client: &http.Client{},
	}
}

func (c *YandexSpeechClient) SynthesizeSpeech(text string, fileName string, accent string) error {
	synURL := "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"
	headers := map[string]string{
		"Authorization": "Api-Key " + c.apiKey,
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	data := url.Values{}
	data.Set("text", text)
	data.Set("lang", "en-En")
	data.Set("format", "mp3")

	req, err := http.NewRequest("POST", synURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		return ioutil.WriteFile(fileName, body, 0644)
	}
	return fmt.Errorf("Error: %s", resp.Status)
}

func (c *YandexSpeechClient) RecognizeSpeech(audioFilePath string) (string, error) {
	audioData, err := ioutil.ReadFile(audioFilePath)
	if err != nil {
		return "", err
	}

	urlRec := "https://stt.api.cloud.yandex.net/speech/v1/stt:recognize?lang=en-US"
	req, err := http.NewRequest("POST", urlRec, bytes.NewBuffer(audioData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Api-Key "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var answer Answer
	err = json.Unmarshal(body, &answer)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 200 {
		return answer.Result, nil
	}
	return "", fmt.Errorf("Error: %s", resp.Status)
}
