package webhook

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	blue   *color.Color
	green  *color.Color
	red    *color.Color
	yellow *color.Color
	gray   *color.Color
}

type ForwardedLog struct {
	evt        WebhookEvent
	statusCode int
	duration   time.Duration
	body       []byte
}

func NewLogger() *Logger {
	return &Logger{
		blue:   color.New(color.FgBlue, color.Bold),
		green:  color.New(color.FgGreen, color.Bold),
		red:    color.New(color.FgRed, color.Bold),
		yellow: color.New(color.FgYellow, color.Bold),
		gray:   color.New(color.FgHiBlack),
	}
}

var timestamp = time.Now().Format("08:06:20")

func (l *Logger) LogReceived(evt WebhookEvent) {
	l.gray.Printf("[%s]", timestamp)
	l.green.Printf("-> ")
	fmt.Printf("Received: %s", evt.Type)
	l.gray.Printf("(ID: %s)\n", evt.ID)
}

func (l *Logger) LogForwarded(fl *ForwardedLog) {
	l.gray.Printf("[%s]", timestamp)
}
