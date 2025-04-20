package logger

import (
	"encoding/json"
	"fmt"
	"github.com/Hivemind-Studio/isi-core/internal/constant/loglevel"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

var loggerInstance *zerolog.Logger

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	// JSON logger for production
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", "isi-core").
		Logger()

	log.Logger = logger
	loggerInstance = &logger
}

func GetLogger() *zerolog.Logger {
	return loggerInstance
}

func Print(logLevel string, requestId string, className string, functionName string, message string, params interface{}) {
	var formattedParams string

	switch v := params.(type) {
	case string:
		formattedParams = v
	case error:
		formattedParams = v.Error()
	default:
		prettyParams, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			formattedParams = fmt.Sprintf("Error marshalling params: %v", err)
		} else {
			formattedParams = string(prettyParams)
		}
	}

	if unescapedStr, err := strconv.Unquote(`"` + formattedParams + `"`); err == nil {
		formattedParams = unescapedStr
	}

	event := log.With().
		Str("request_id", requestId).
		Str("class", className).
		Str("function", functionName).
		Str("message", message).
		Str("parameters", formattedParams).
		Logger()

	switch logLevel {
	case loglevel.INFO:
		event.Info().Msg(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))
	case loglevel.WARN:
		event.Warn().Msg(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))
	case loglevel.DEBUG:
		event.Debug().Msg(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))
	case loglevel.ERROR:
		event.Error().Msg(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))
	default:
		event.Info().Msg(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))
	}
}
