package authorization

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
)

type Authorization struct {
	logger *zap.Logger
}

func NewAuthorization(logger *zap.Logger) (*Authorization, error){
	return &Authorization{
		logger: logger,
	}, nil
}

func (svc *Authorization) CheckCredentials(reqId string, username string, password string) bool{
	svc.logger.Debug("Checking credentials",
		zap.String("reqId", reqId),
		zap.String("user", username),
	)
	if pass,found := os.LookupEnv(fmt.Sprintf("USER_%s", strings.ToUpper(username))); found {
		if pass == password {
			return true
		}
	}
	return false
}
