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
	conn          *websocket.Conn
	httpClient    *http.Client
	logger        *Logger
}

type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	CreatedAt time.Time              `json:"created_at"`
	Data      map[string]interface{} `json:"data"`
}

var websocketURL = "wss://api.abacatepay.com/v1/webhooks/stream"

func genSigningSecret() string {
	return fmt.Sprintf("whsec_%d", time.Now().UnixNano())
}

func NewProxy(apiKey, forwardTo string) *Proxy {
	return &Proxy{
		apiKey:        apiKey,
		forwardTo:     forwardTo,
		signingSecret: genSigningSecret(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *Proxy) Connect() error {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+p.apiKey)
	header.Set("User-Agent", "Abacate-CLI/1.0")

	conn, reply, err := websocket.DefaultDialer.Dial(websocketURL, header)
	if err != nil {
		if reply != nil {
			body, _ := io.ReadAll(reply.Body)
			return fmt.Errorf("error to connect: (status %d) %s", reply.StatusCode, string(body))
		}
		return err
	}

	p.conn = conn
	return nil
}

func (p *Proxy) Listen() error {
	defer p.conn.Close()

	for {
		var evt WebhookEvent

		err := p.conn.ReadJSON(&evt)
		if err != nil {

			// Default websocket closing error
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}

			return fmt.Errorf("error to read event: %w", err)
		}

		go p.handleEvent(evt)
	}
}

func (p *Proxy) handleEvent(evt WebhookEvent) {
	start := time.Now()

	p.logger.LogReceived(evt)

	payload, err := json.Marshal(evt)
	if err != nil {
		p.logger.LogError(evt, fmt.Errorf("error to serialize payload: %w", err))
	}

	signature := p.generateSignature(payload, evt.CreatedAt)
	req, err := http.NewRequest("POST", "http://re"+p.forwardTo, bytes.NewBuffer(payload))
	if err != nil {
		p.logger.LogError(evt, fmt.Errorf("error to create request: %w", err))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Abacate-Signature", signature)
	req.Header.Set("X-Abacate-Event-Type", evt.Type)
	req.Header.Set("X-Abacate-Event-ID", evt.ID)

	reply, err := p.httpClient.Do(req)
	if err != nil {
		p.logger.LogError(evt, fmt.Errorf("error to make forward: %w", err))
		return
	}
	defer reply.Body.Close()

	body, _ := io.ReadAll(reply.Body)
	duration := time.Since(start)

	p.logger.LogForwarded(&ForwardedLog{
		evt:        evt,
		statusCode: reply.StatusCode,
		body:       body,
		duration:   duration,
	})
}

func (p *Proxy) generateSignature(payload []byte, timestamp time.Time) string {
	signedPayload := fmt.Sprintf("%d,%s", timestamp.Unix(), string(payload))

	h := hmac.New(sha256.New, []byte(p.signingSecret))
	h.Write([]byte(signedPayload))

	return hex.EncodeToString(h.Sum(nil))
}

// TODO: func de Close() pra fechar a conexao
