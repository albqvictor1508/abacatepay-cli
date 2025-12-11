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
	statusCode int
	body       []byte
	event      WebhookEvent
	duration   time.Duration
}

func NewLogger() *Logger {
	return &Logger{
		gray:   color.New(color.FgHiBlack),
		red:    color.New(color.FgRed, color.Bold),
		blue:   color.New(color.FgBlue, color.Bold),
		green:  color.New(color.FgGreen, color.Bold),
		yellow: color.New(color.FgYellow, color.Bold),
	}
}

func Truncate(str string, maxSize int) string {
	if len(str) <= maxSize {
		return str
	}

	return str[:maxSize] + "..."
}

var timestamp = time.Now().Format("08:06:20")

func (logger *Logger) LogReceived(evt WebhookEvent) {
	logger.gray.Printf("[%s]", timestamp)
	logger.green.Printf("-> ")

	fmt.Printf("Received: %s", evt.Type)
	
	logger.gray.Printf("(ID: %s)\n", evt.ID)
}

func (logger *Logger) LogForwarded(fl *ForwardedLog) {
	logger.gray.Printf("[%s]", timestamp)

	logger.PrintStatusIcon(fl.statusCode)

	fmt.Printf("%s", fl.event.Type)

	logger.PrintStatusCode(fl.statusCode)

	logger.gray.Printf("%dms\n", fl.duration.Milliseconds())

	if fl.statusCode >= 400 && len(fl.body) > 0 {
		logger.gray.Printf("Reply:%s\n", Truncate(string(fl.body), 100))
	}
}

func (logger *Logger) LogError(event WebhookEvent, err error) {
	logger.gray.Printf("[%s]", timestamp)

	logger.red.Printf("✗ ")

	fmt.Printf("Error processing event %s (ID: %s): %v\n", event.Type, event.ID, err)
}

func (logger *Logger) PrintStatusIcon(status int) {
	if status >= 400 {
		logger.red.Printf("✗ ")
		logger.red.Printf("[%d]", status)

		return
	}

	if status >= 200 && status < 300 {
		logger.green.Printf("✓ ")
		logger.green.Printf("[%d]", status)

		return
	}

	logger.yellow.Printf("→ ")
	logger.yellow.Printf("[%d]", status)
}

func (logger *Logger) PrintStatusCode(status int) {
	if status >= 400 {
		logger.red.Printf("[%d]", status)

		return
	}

	if status >= 200 && status < 300 {
		logger.green.Printf("[%d]", status)

		return
	}

	logger.yellow.Printf("[%d]", status)
}
