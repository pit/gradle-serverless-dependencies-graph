package authorization

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
	"terraform-serverless-private-registry/lib/helpers"
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
	envVarName := fmt.Sprintf("USER_%s", strings.ToUpper(helpers.GenerateMD5(username)))
	svc.logger.Debug("Looking to env var",
		zap.String("env var name", envVarName),
	)
	if pass,found := os.LookupEnv(envVarName); found {
		svc.logger.Debug("Env var found",
			zap.String("env var name", envVarName),
			zap.String("env var value", pass),
			)
		if pass == helpers.GenerateMD5(password) {
			return true
		}
	}
	return false
}

