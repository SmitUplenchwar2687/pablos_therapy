package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func generateAndPlayAudio(text string) {
	data := map[string]interface{}{
		"text":  text,
		"voice": "George",
		"model": "eleven_multilingual_v2",
	}
	payload, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", "https://api.elevenlabs.io/v1/text-to-speech", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", elevenlabsAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}

	defer resp.Body.Close()
	audio, _ := io.ReadAll(resp.Body)
	playAudio(audio)
}
