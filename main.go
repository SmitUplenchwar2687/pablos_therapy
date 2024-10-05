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

	fmt.Println("Starting GO API service...")

	fmt.Println(`
 ______     ______        ______     ______   __    
/\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
\ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
 \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
  \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)

	r.Get("/voice-to-text", func(w http.ResponseWriter, r *http.Request) {
		text := voiceToText()
		w.Write([]byte(text))
	})

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
	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log1.Error(err)
	}

	text := voiceToText()
	response := chatWithGPT(text)
	generateAndPlayAudio(response)

}
