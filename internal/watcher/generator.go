package watcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Generator struct {
	apiKey      string
	model       string
	temperature float64
}

func NewGenerator(apiKey, model string, temperature float64) *Generator {
	return &Generator{apiKey: apiKey, model: model, temperature: temperature}
}

func (g *Generator) Generate(prompt string) (string, error) {
	url := "https://api.openai.com/v1/completions"
	body := map[string]interface{}{
		"model":       g.model,
		"prompt":      prompt,
		"temperature": g.temperature,
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("failed to call OpenAI API")
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["choices"].([]interface{})[0].(map[string]interface{})["text"].(string), nil
}
