package lib

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

type Modules struct {
	storage *Storage
	logger  *zap.Logger
}

type ListModuleVersionsResponse struct {
	Modules []ModuleVersions `json:"modules"`
}

type ModuleVersions struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Version string `json:"version"`
}

func NewModules(storageObj *Storage, loggerObj *zap.Logger) (*Modules, error) {
	return &Modules{
		storage: storageObj,
		logger: loggerObj,
	}, nil
}

func (svc *Modules) ListModuleVersions(ctx context.Context, namespace string, name string, provider string) (ListModuleVersionsResponse, error) {
	svc.logger.Debug("modules.ListModuleVersions() called",
		zap.String("namespace", namespace),
		zap.String("name", name),
		zap.String("provider", provider),
		)
	dirPath := fmt.Sprintf("/modules/%s/%s/%s", namespace, name, provider)
	dirs,_ := svc.storage.ListDirs(dirPath)

	var versions []Version
	for _,dir := range dirs{
		versions = append(versions, Version{Version: dir})
	}

	result := ListModuleVersionsResponse{
		Modules: []ModuleVersions{{Versions : versions}},
	}

	svc.logger.Debug("ListModuleVersions() return",
		zap.Reflect("result", result),
		zap.String("namespace", namespace),
		zap.String("name", name),
		zap.String("provider", provider),
	)
	return result, nil
}
