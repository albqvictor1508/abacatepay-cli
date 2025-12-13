package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func TriggerLocalEvent(eventType, forwardTo, signingSecret string) error {
	templates := GetEventTemplates()
	template, exists := templates[eventType]

	if !exists {
		return fmt.Errorf("invalid event type: %s", eventType)
	}

	event := WebhookEvent{
		ID:        "evt_" + uuid.New().String()[:12],
		Type:      template.Type,
		CreatedAt: time.Now(),
		Data:      template.Data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error to serialize json: %w", err)
	}

	signature := generateSignature(payload, event.CreatedAt, signingSecret)

	req, err := http.NewRequest("POST", "http://"+forwardTo, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Abacate-Signature", signature)
	req.Header.Set("X-Abacate-Event-Type", event.Type)
	req.Header.Set("X-Abacate-Event-ID", event.ID)
	req.Header.Set("X-Abacate-Test-Event", "true")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("endpoint returns status %d", resp.StatusCode)
	}

	return nil
}
