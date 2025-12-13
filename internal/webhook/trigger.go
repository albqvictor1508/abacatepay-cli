package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dimiro1/faker"
	"github.com/google/uuid"
)

type EventLog struct {
	Event     *WebhookEvent `json:"event"`
	Timestamp time.Time     `json:"timestamp"`
	Source    string        `json:"source"`
}

type EventTemplate struct {
	Type string
	Data map[string]interface{}
}

func GetEventTemplates() map[string]EventTemplate {
	customerPF := map[string]interface{}{
		"id":       "cust_" + uuid.New().String()[:8],
		"name":     faker.Name(),
		"email":    faker.Internet().Email(),
		"phone":    faker.PhoneNumber().CellPhone(),
		"document": faker.Business().Cpf(),
		"type":     "individual",
	}

	customerPJ := map[string]interface{}{
		"id":           "cust_" + uuid.New().String()[:8],
		"name":         faker.Company().Name(),
		"email":        faker.Internet().FreeEmail(),
		"phone":        faker.PhoneNumber().PhoneNumber(),
		"document":     faker.Business().Cnpj(),
		"type":         "business",
		"company_name": faker.Company().Name(),
	}

	customer := customerPF
	if rand.Intn(2) == 0 {
		customer = customerPJ
	}

	amounts := []int{1990, 2990, 4990, 7990, 9990, 14990, 19990, 29990, 49990, 99990}
	amount := amounts[rand.Intn(len(amounts))]

	products := []string{
		"Plano Básico",
		"Plano Premium",
		"Plano Enterprise",
		"Plano Starter",
		"Plano Pro",
	}
	selectedProduct := products[rand.Intn(len(products))]

	return map[string]EventTemplate{
		"billing.paid": {
			Type: "billing.paid",
			Data: map[string]interface{}{
				"id":       "bill_" + uuid.New().String()[:8],
				"amount":   amount,
				"currency": "BRL",
				"status":   "paid",
				"paid_at":  time.Now().Format(time.RFC3339),
				"pix": map[string]interface{}{
					"txid":       faker.Lorem().Characters(32),
					"end_to_end": "E" + faker.Lorem().Characters(31),
					"payer": map[string]interface{}{
						"name":     customer["name"],
						"document": customer["document"],
					},
				},
				"customer": customer,
				"product": map[string]interface{}{
					"id":          "prod_" + uuid.New().String()[:8],
					"name":        selectedProduct,
					"description": "Assinatura " + selectedProduct,
				},
				"metadata": map[string]interface{}{
					"order_id":   faker.Lorem().Characters(16),
					"ip_address": faker.Internet().IpV4Address(),
					"user_agent": faker.Internet().UserAgent(),
				},
				"created_at": time.Now().Add(-2 * time.Minute).Format(time.RFC3339),
			},
		},
		"withdraw.done": {
			Type: "withdraw.done",
			Data: map[string]interface{}{
				"id":           "withdraw_" + uuid.New().String()[:8],
				"amount":       amount,
				"currency":     "BRL",
				"status":       "completed",
				"completed_at": time.Now().Format(time.RFC3339),
				"method":       "pix",
				"pix": map[string]interface{}{
					"txid":       faker.Lorem().Characters(32),
					"end_to_end": "E" + faker.Lorem().Characters(31),
					"recipient": map[string]interface{}{
						"name":         faker.Name(),
						"document":     faker.Business().Cpf(),
						"bank":         []string{"Nubank", "Banco do Brasil", "Itaú", "Bradesco", "Santander", "Caixa"}[rand.Intn(6)],
						"pix_key":      faker.Internet().Email(),
						"pix_key_type": "email",
					},
				},
				"fee": map[string]interface{}{
					"amount":      int(float64(amount) * 0.02),
					"currency":    "BRL",
					"description": "Taxa de saque",
				},
				"net_amount": int(float64(amount) * 0.98),
				"metadata": map[string]interface{}{
					"requested_by": faker.Internet().Email(),
					"ip_address":   faker.Internet().IpV4Address(),
				},
				"requested_at": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			},
		},
		"withdraw.failed": {
			Type: "withdraw.failed",
			Data: map[string]interface{}{
				"id":        "withdraw_" + uuid.New().String()[:8],
				"amount":    amount,
				"currency":  "BRL",
				"status":    "failed",
				"failed_at": time.Now().Format(time.RFC3339),
				"method":    "pix",
				"failure_reason": []string{
					"insufficient_balance",
					"invalid_pix_key",
					"recipient_account_blocked",
					"daily_limit_exceeded",
					"pix_key_not_found",
					"technical_error",
				}[rand.Intn(6)],
				"failure_message": map[string]string{
					"insufficient_balance":      "Saldo insuficiente para realizar o saque",
					"invalid_pix_key":           "Chave PIX inválida ou não encontrada",
					"recipient_account_blocked": "Conta do destinatário bloqueada",
					"daily_limit_exceeded":      "Limite diário de saques excedido",
					"pix_key_not_found":         "Chave PIX não cadastrada",
					"technical_error":           "Erro técnico ao processar saque",
				}[[]string{
					"insufficient_balance",
					"invalid_pix_key",
					"recipient_account_blocked",
					"daily_limit_exceeded",
					"pix_key_not_found",
					"technical_error",
				}[rand.Intn(6)]],
				"pix": map[string]interface{}{
					"recipient": map[string]interface{}{
						"name":         faker.Name(),
						"document":     faker.Business().Cpf(),
						"bank":         []string{"Nubank", "Banco do Brasil", "Itaú", "Bradesco", "Santander"}[rand.Intn(5)],
						"pix_key":      faker.Internet().Email(),
						"pix_key_type": "email",
					},
				},
				"metadata": map[string]interface{}{
					"requested_by": faker.Internet().Email(),
					"ip_address":   faker.Internet().IpV4Address(),
					"retry_count":  rand.Intn(3),
				},
				"requested_at": time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
			},
		},
	}
}

func generateSignature(payload []byte, timestamp time.Time, secret string) string {
	signedPayload := fmt.Sprintf("%d.%s", timestamp.Unix(), string(payload))
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signedPayload))
	return hex.EncodeToString(h.Sum(nil))
}

func ListAvailableEvents() []string {
	templates := GetEventTemplates()
	events := make([]string, 0, len(templates))
	for eventType := range templates {
		events = append(events, eventType)
	}
	return events
}

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
