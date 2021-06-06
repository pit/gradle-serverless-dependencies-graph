package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strings"
)

type ResponseNotFound struct {
	Status *string `json:"status"`
}

func ApiNotFound() *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusNotFound

	status := "Not found"
	stringBody, err := json.Marshal(&ResponseNotFound{
		Status: &status,
	})
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response json: %s", string(stringBody))

	return &resp
}
func ApiNoContent() *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusNoContent

	return &resp
}

func ApiResponse(status int, body interface{}) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = status

	stringBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response json: %s", string(stringBody))

	return &resp
}

func InitLogger(logLevel string, structured bool) (*zap.Logger, error) {
	spew.Config.Indent = "  "
	spew.Config.DisableMethods = true
	spew.Config.DisablePointerMethods = true

	var lvl zap.AtomicLevel
	switch strings.ToLower(logLevel) {
	case "info":
		lvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn", "warning":
		lvl = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error", "err":
		lvl = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		lvl = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	encoding := "console"
	if structured {
		encoding = "json"
	}

	cfg := zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         encoding,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()

	return logger, err

}
