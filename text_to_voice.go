package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gordonklaus/portaudio"
)

const (
	channels        = 1
	framesPerBuffer = 1024
	sampleRate      = 44100
)

func GenerateAudio(text string) ([]byte, error) {
	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_multilingual_v2",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.elevenlabs.io/v1/text-to-speech/JBFqnCBsd6RMkjVDRZzb", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", elevenlabsAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response status: %s, body: %s", resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}

func convertMp3ToPcm(mp3Data []byte) ([]float32, error) {
	tempDir, err := os.MkdirTemp("", "audio")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write MP3 data to temporary file
	mp3Path := filepath.Join(tempDir, "audio.mp3")
	if err := os.WriteFile(mp3Path, mp3Data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write MP3 file: %v", err)
	}

	// Convert MP3 to 16-bit PCM using ffmpeg
	rawPath := filepath.Join(tempDir, "audio.raw")
	cmd := exec.Command("ffmpeg",
		"-i", mp3Path,
		"-f", "s16le", // 16-bit signed integer PCM
		"-acodec", "pcm_s16le",
		"-ar", fmt.Sprintf("%d", sampleRate),
		"-ac", "1", // mono
		rawPath)

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg conversion failed: %v\nOutput: %s", err, string(output))
	}

	rawData, err := os.ReadFile(rawPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PCM file: %v", err)
	}

	numSamples := len(rawData) / 2 // 2 bytes per int16
	samples := make([]float32, numSamples)

	for i := 0; i < numSamples; i++ {
		sample := int16(binary.LittleEndian.Uint16(rawData[i*2 : (i+1)*2]))
		samples[i] = float32(sample) / 32768.0
	}

	fmt.Printf("Converted audio: Samples: %d\n", len(samples))

	maxAmp := float32(0)
	for _, sample := range samples {
		amp := abs(sample)
		if amp > maxAmp {
			maxAmp = amp
		}
	}

	if maxAmp > 1.0 {
		for i := range samples {
			samples[i] /= maxAmp
		}
	} else if maxAmp < 0.1 {
		gain := 0.5 / maxAmp
		for i := range samples {
			samples[i] *= gain
		}
	}

	fmt.Printf("Audio max amplitude after normalization: %f\n", maxAmp)

	return samples, nil
}

func PlayAudio(audioData []byte) error {
	samples, err := convertMp3ToPcm(audioData)
	if err != nil {
		return fmt.Errorf("failed to convert audio: %v", err)
	}

	if len(samples) == 0 {
		return fmt.Errorf("no audio samples to play")
	}

	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize portaudio: %v", err)
	}
	defer portaudio.Terminate()

	buffer := make([]float32, framesPerBuffer)
	stream, err := portaudio.OpenDefaultStream(0, channels, float64(sampleRate), framesPerBuffer, buffer)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %v", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return fmt.Errorf("failed to start audio stream: %v", err)
	}

	fmt.Println("Playing audio...")

	for i := 0; i < len(samples); i += framesPerBuffer {
		end := i + framesPerBuffer
		if end > len(samples) {
			end = len(samples)
			for j := end - i; j < framesPerBuffer; j++ {
				buffer[j] = 0
			}
		}

		copy(buffer, samples[i:end])

		if err := stream.Write(); err != nil {
			return fmt.Errorf("failed to write to audio stream: %v", err)
		}
	}

	if err := stream.Stop(); err != nil {
		return fmt.Errorf("failed to stop audio stream: %v", err)
	}

	return nil
}

func abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}
