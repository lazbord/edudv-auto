package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Payload struct {
	Content string `json:"content"`
}

func PostReqToDiscord(webhookURL, message string) error {
	payload := Payload{
		Content: message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response from Discord: %v", resp.Status)
	}

	return nil
}

func SendDiscordMessage(message string) {
	godotenv.Load()

	webhookURL := os.Getenv("webhookURL")

	err := PostReqToDiscord(webhookURL, message)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}
}

func SendLoggerMessage(message string) {
	godotenv.Load()

	webhookURL := os.Getenv("loggerwebhookURL")

	err := PostReqToDiscord(webhookURL, message)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}
}
