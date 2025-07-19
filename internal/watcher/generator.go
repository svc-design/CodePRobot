package watcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
)

type Generator struct {
	apiKey      string
	model       string
	temperature float64
	apiURL      string
}

func NewGenerator(apiKey, model, apiURL string, temperature float64) *Generator {
	if apiURL == "" {
		apiURL = "https://api.openai.com/v1/completions"
	}
	return &Generator{apiKey: apiKey, model: model, apiURL: apiURL, temperature: temperature}
}

func (g *Generator) Generate(prompt string) (string, error) {
	url := g.apiURL
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

// GenerateWithCodex 调用 codex-cli 生成代码
func (g *Generator) GenerateWithCodex(prompt string) (string, error) {
	cmd := exec.Command("npx", "codex-cli", "--prompt", prompt)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// GenerateWithClaude 调用 claude-code 包生成代码
func (g *Generator) GenerateWithClaude(prompt string) (string, error) {
	cmd := exec.Command("npx", "claude-code", prompt)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
