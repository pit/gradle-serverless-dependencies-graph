package helpers

import (
	"crypto/md5"
	"encoding/hex"
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

func ApiErrorNotFound() *events.APIGatewayProxyResponse {
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


func ApiErrorUnknown() *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusInternalServerError

	status := "Unknown"
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

func ApiErrorNoContent() *events.APIGatewayProxyResponse {
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

	fmt.Printf("response: status=%d, json: %s", resp.StatusCode, string(stringBody))

	return &resp
}

func HtmlResponse(status int, body *string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "text/html"},
	}
	resp.StatusCode = status

	resp.Body = *body

	fmt.Printf("response: status=%d, body: %s", resp.StatusCode, *body)

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

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func GenerateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GenerateMD5List(texts []string) []string {
	result := make([]string, len(texts))
	for idx, val := range texts {
		result[idx] = GenerateMD5(val)
	}
	return result
}
