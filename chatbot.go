package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func chatWithGPT(prompt string) string {
	data := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a multilingual therapy chatbot."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.9,
	}
	payload, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	fmt.Printf("GPT Response: %s\n", message)
	return message
}
