// voice.go
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

// VoiceToText function to handle audio recording and conversion
func VoiceToText() string {
	audioFile := "output.wav"
	err := recordAudio(audioFile)
	if err != nil {
		fmt.Printf("Error recording audio: %v\n", err)
		return ""
	}

	text, err := convertAudioToText(audioFile)
	if err != nil {
		fmt.Printf("Error converting audio to text: %v\n", err)
		return ""
	}
	return text
}

// Record audio using arecord (Linux) or equivalent
func recordAudio(outputFile string) error {
	fmt.Println("Recording audio for 2 seconds...")
	cmd := exec.Command("sox", "-d", "-c", "1", "-r", "16000", "-b", "16", outputFile, "trim", "0", "2")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to record audio: %v", err)
	}
	fmt.Println("Recording completed.")
	return nil
}

// Convert the recorded audio file to text
func convertAudioToText(audioFile string) (string, error) {
	// Create a new Google Speech-to-Text client
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Read the audio file into memory
	audioData, err := ioutil.ReadFile(audioFile)
	if err != nil {
		return "", fmt.Errorf("failed to read audio file: %v", err)
	}

	// Configure the request with audio content and recognition settings
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000, // Adjust based on your audio file
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioData},
		},
	}

	// Perform the speech recognition request
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to recognize speech: %v", err)
	}

	// Check if there are results
	if len(resp.Results) == 0 {
		return "", fmt.Errorf("no transcription result found")
	}

	// Concatenate all transcriptions from each result segment
	var transcription string
	for _, result := range resp.Results {
		transcription += result.Alternatives[0].Transcript + " "
	}

	return transcription, nil
}
