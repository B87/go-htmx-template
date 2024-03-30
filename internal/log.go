package internal

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
)

type Logger interface {
	// Log logs a message
	Log(message string)
	// Error logs an error message
	Error(message string)
}
