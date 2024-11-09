package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	log1 "github.com/sirupsen/logrus"
)

const openaiAPIKey = ""
const elevenlabsAPIKey = ""

func main() {

	log1.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()

	// Get the transcribed text from audio
	transcribedText := VoiceToText()
	fmt.Println("Transcribed Text:", transcribedText)

	// Pass the transcribed text to chatWithGPT and get the response
	gptResponse := chatWithGPT(transcribedText)
	fmt.Println("GPT-3 Response:", gptResponse)

	// Pass the transcribed text to chatWithGPT (assuming chatWithGPT is defined in main.go)
	// if transcribedText != "" {
	// 	gptResponse := chatWithGPT(transcribedText)
	// 	fmt.Println("GPT-3 Response:", gptResponse)
	// } else {
	// 	fmt.Println("No transcribed text available.")
	// }

	// Generate audio from the text using GenerateAudio function
	audioData, err := GenerateAudio(gptResponse)
	if err != nil {
		fmt.Printf("Error generating audio: %v\n", err)
		return
	}

	// Play the generated audio using PlayAudio function
	if err := PlayAudio(audioData); err != nil {
		fmt.Printf("Error playing audio: %v\n", err)
	}

	fmt.Println("Starting GO API service...")

	fmt.Println(`
 ______     ______        ______     ______   __    
/\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
\ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
 \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
  \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)

	// HTTP endpoint for voice-to-text
	r.Get("/voice-to-text", func(w http.ResponseWriter, r *http.Request) {
		text := VoiceToText()
		fmt.Println("Transcribed Text:", text)
		w.Write([]byte(text))
	})

	// HTTP endpoint to generate response using chatWithGPT
	r.Post("/generate-response", func(w http.ResponseWriter, r *http.Request) {
		var prompt struct {
			Text string `json:"text"`
		}
		err := json.NewDecoder(r.Body).Decode(&prompt)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		response := chatWithGPT(prompt.Text)
		w.Write([]byte(response))
	})

	// Start the HTTP server
	// err := http.ListenAndServe("localhost:8000", r)
	// if err != nil {
	// 	log1.Error(err)
	// }
}
