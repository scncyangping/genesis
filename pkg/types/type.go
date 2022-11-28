package types

import (
	"path/filepath"
)

const SP = string(filepath.Separator)

type B map[string]any

type LogLevel int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel LogLevel = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

var LogLevelM = map[string]LogLevel{
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
}

type AppMod string

const (
	// DebugMode indicates sg mode is debug.
	DebugMode AppMod = "debug"
	// ReleaseMode indicates sg mode is release.
	ReleaseMode AppMod = "release"
)
