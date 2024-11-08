package logger

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

func Print(logLevel string, requestId string, className string, functionName string, message string, params interface{}) {
	var log *logrus.Entry
	fields := logrus.Fields{
		"request_id": requestId,
		"class":      className,
		"function":   functionName,
		"message":    message,
	}

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

	unescapedStr, err := strconv.Unquote(`"` + formattedParams + `"`)
	if err == nil {
		formattedParams = unescapedStr
	}

	switch logLevel {
	case "info":
		fields["parameters"] = formattedParams
		log = logrus.WithFields(fields)
		log.Info(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))

	case "warn":
		fields["parameters"] = formattedParams
		log = logrus.WithFields(fields)
		log.Warn(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))

	case "debug":
		fields["parameters"] = formattedParams
		log = logrus.WithFields(fields)
		log.Debug(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))

	case "error":
		fields["parameters"] = formattedParams
		log = logrus.WithFields(fields)
		log.Error(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))

	default:
		log = logrus.WithFields(fields)
		log.Info(fmt.Sprintf("Function: %s processed by class: %s", functionName, className))
	}
}
