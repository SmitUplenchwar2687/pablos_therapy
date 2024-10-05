package main

import (
	"github.com/gordonklaus/portaudio"
	"log"
)

func playAudio(audio []byte) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(audio), &audio)
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}
	defer stream.Close()

	stream.Start()
	defer stream.Stop()

	stream.Write()
}
