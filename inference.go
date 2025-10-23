package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type llm struct {
	apiUrl    string
	modelCard string
	provider  string
	apiKey    string
}

var providerMap map[string]string = map[string]string{
	"groq": "https://api.groq.com/openai/v1/chat/completions",
}

func InitialiseModel(modelCard, provider string) llm {
	var newLlm llm

	newLlm.apiUrl = providerMap[provider]
	newLlm.modelCard = modelCard
	newLlm.provider = provider
	newLlm.apiKey = os.Getenv("GROQ_API_KEY")

	return newLlm
}
func (m llm) invoke(prompt string) (string, error) {
	client := &http.Client{}
	payload := map[string]interface{}{
		"model": m.modelCard,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 1,
	}
	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", m.apiUrl, bytes.NewBuffer(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+m.apiKey) // no trailing space

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("non-200 status: %d, body: %s", resp.StatusCode, body)
	}

	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return "", fmt.Errorf("json unmarshal error: %w", err)
	}

	choices := v["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"]
	content := message.(map[string]interface{})["content"].(string)

	return content, nil
}
