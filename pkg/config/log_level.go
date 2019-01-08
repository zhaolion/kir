package config

type LogLevel int

// LogLevel set
const (
	LogDebug LogLevel = iota
	LogInfo
	LogWarning
	LogError
	LogOff
)
