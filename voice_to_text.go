package main

import (
	"context"
	"fmt"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
)

func voiceToText() string {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{
				Uri: "gs://your-bucket/your-audio-file.wav",
			},
		},
	}
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		log.Fatalf("Failed to recognize: %v", err)
	}
	var resultText string
	for _, result := range resp.Results {
		resultText += result.Alternatives[0].Transcript
	}
	fmt.Printf("Recognized text: %s\n", resultText)

	return resultText
}
