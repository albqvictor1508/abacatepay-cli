package webhook

import (
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

	p.logger.Log

	json.Unmarshal(evt)
}
