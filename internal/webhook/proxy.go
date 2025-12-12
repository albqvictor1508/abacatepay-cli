package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Proxy struct {
	apiKey        string
	forwardTo     string
	signingSecret string
	logger        *Logger
	client        *http.Client
	connection    *websocket.Conn
}

type WebhookEvent struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Data      map[string]any `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
}

const websocketURL = "wss://api.abacatepay.com/v1/webhooks/stream"

func genSigningSecret() string {
	return fmt.Sprintf("whsec_%d", time.Now().UnixNano())
}

func (p *Proxy) getSigningSecret() string {
	return p.signingSecret
}

func NewProxy(apiKey, forwardTo string) *Proxy {
	return &Proxy{
		apiKey:        apiKey,
		forwardTo:     forwardTo,
		signingSecret: genSigningSecret(),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (proxy *Proxy) Connect() error {
	header := http.Header{}

	header.Set("User-Agent", "AbacatePayCLI (1.0.0)")
	header.Set("Authorization", "Bearer "+proxy.apiKey)

	connection, reply, err := websocket.DefaultDialer.Dial(websocketURL, header)
	if err != nil {
		if reply != nil {
			body, _ := io.ReadAll(reply.Body)

			return fmt.Errorf("connection error (%d): %s", reply.StatusCode, string(body))
		}

		return err
	}

	proxy.connection = connection

	return nil
}

func (proxy *Proxy) Listen() error {
	defer proxy.connection.Close()

	for {
		var evt WebhookEvent

		if err := proxy.connection.ReadJSON(&evt); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}

			return fmt.Errorf("unknown error while reading event: %w", err)
		}

		go proxy.Handle(evt)
	}
}

func (proxy *Proxy) Handle(event WebhookEvent) {
	start := time.Now()

	proxy.logger.LogReceived(event)

	payload, err := json.Marshal(event)
	if err != nil {
		proxy.logger.LogError(event, fmt.Errorf("error to serialize payload: %w", err))
	}

	signature := proxy.GenerateSignature(payload, event.CreatedAt)
	req, err := http.NewRequest("POST", "http://re"+proxy.forwardTo, bytes.NewBuffer(payload))
	if err != nil {
		proxy.logger.LogError(event, fmt.Errorf("error to create request: %w", err))

		return
	}

	req.Header.Set("X-Abacate-Event-ID", event.ID)
	req.Header.Set("X-Abacate-Signature", signature)
	req.Header.Set("X-Abacate-Event-Type", event.Type)
	req.Header.Set("Content-Type", "application/json")

	reply, err := proxy.client.Do(req)
	if err != nil {
		proxy.logger.LogError(event, fmt.Errorf("forward Error: %w", err))

		return
	}

	defer reply.Body.Close()

	duration := time.Since(start)
	body, err := io.ReadAll(reply.Body)
	if err != nil {
		proxy.logger.LogError(event, fmt.Errorf("body Download Error: %w", err))
	}

	proxy.logger.LogForwarded(&ForwardedLog{
		event:      event,
		body:       body,
		duration:   duration,
		statusCode: reply.StatusCode,
	})
}

func (proxy *Proxy) GenerateSignature(payload []byte, timestamp time.Time) string {
	signed := fmt.Sprintf("%d,%s", timestamp.Unix(), string(payload))

	hash := hmac.New(sha256.New, []byte(proxy.signingSecret))

	hash.Write([]byte(signed))

	return hex.EncodeToString(hash.Sum(nil))
}

func (p *Proxy) Close() {
	if p.connection != nil {
		p.connection.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}

	p.connection.Close()
}
