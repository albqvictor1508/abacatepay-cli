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

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
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

	l.printStatusIcon(fl.statusCode)
	fmt.Printf("%s", fl.evt.Type)

	l.printStatusCode(fl.statusCode)

	l.gray.Printf("%dms\n", fl.duration.Milliseconds())
	if fl.statusCode >= 400 && len(fl.body) > 0 {
		l.gray.Printf("Reply:%s\n", truncate(string(fl.body), 100))
	}
}

func (l *Logger) printStatusIcon(statusCode int) {
	if statusCode >= 400 {
		l.red.Printf("✗ ")
		l.red.Printf("[%d]", statusCode)
		return
	}

	if statusCode >= 200 && statusCode < 300 {
		l.green.Printf("✓ ")
		l.green.Printf("[%d]", statusCode)
		return
	}

	l.yellow.Printf("→ ")
	l.yellow.Printf("[%d]", statusCode)
}

func (l *Logger) printStatusCode(statusCode int) {
	if statusCode >= 400 {
		l.red.Printf("[%d]", statusCode)
		return
	}

	if statusCode >= 200 && statusCode < 300 {
		l.green.Printf("[%d]", statusCode)
		return
	}

	l.yellow.Printf("[%d]", statusCode)
}
