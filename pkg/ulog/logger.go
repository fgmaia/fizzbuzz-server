package ulog

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	logger      zerolog.Logger
	gitRevision string
	buildInfo   *debug.BuildInfo
)

func LogInit() {
	// Console Writer (for stdout)
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	// Combine Writers (Console + File)
	multiWriter := zerolog.MultiLevelWriter(consoleWriter)

	// Create the logger
	logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	// Set global log level (optional)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// default fields, remover if not needed
	// Get build info
	var ok bool
	if buildInfo, ok = debug.ReadBuildInfo(); ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
	}

	// Add program information as a context group
	//logger = logger.With().
	//	Int("pid", os.Getpid()).
	//	Str("go_version", buildInfo.GoVersion).
	//	Str("vcs_revision", gitRevision).
	//	Logger()
}

func Info(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		var sb strings.Builder
		for _, field := range fields {
			sb.WriteString(fmt.Sprintf("%v ", field))
		}
		logger.Info().Msg(msg + " " + sb.String())
	} else {
		logger.Info().Msg(msg)
	}
}

func Infof(msg string, fields ...interface{}) {
	logger.Info().Msgf(msg, fields...)
}

func InfoE(msg string) {
	logger.Info().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Str("vcs_revision", gitRevision).
		Msg(msg)
}

func Error(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		var sb strings.Builder
		for _, field := range fields {
			sb.WriteString(fmt.Sprintf("%v ", field))
		}
		logger.Error().Msg(msg + " " + sb.String())
	} else {
		logger.Error().Msg(msg)
	}
}

func Errorf(msg string, fields ...interface{}) {
	logger.Error().Msgf(msg, fields...)
}

func Panic(r any) {
	logger.Error().
		Msg(string(debug.Stack()))
	logger.Error().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Str("vcs_revision", gitRevision).
		Msg("Application crashed")
}
