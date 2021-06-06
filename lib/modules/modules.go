package modules

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"terraform-serverless-private-registry/lib/storage"
)

type Modules struct {
	storageSvc *storage.Storage
	logger     *zap.Logger
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

type InputParams struct {
	Namespace *string
	Name      *string
	Provider  *string
	Version   *string
}

const (
	ErrNotFound = iota
	ErrUnknown
)

type ModulesError struct {
	Message string
	Code    int
	Err     error
}

func (s ModulesError) Error() string {
	panic("implement me")
}

func NewModules(storage *storage.Storage, log *zap.Logger) (*Modules, error) {
	return &Modules{
		storageSvc: storage,
		logger:     log,
	}, nil
}

func (svc *Modules) ListModuleVersions(ctxId string, params InputParams) (*ListModuleVersionsResponse, *ModulesError) {
	svc.logger.Debug(fmt.Sprintf("%s modules.ListModuleVersions() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Name),
		zap.String("provider", *params.Provider),
	)
	dirPath := fmt.Sprintf("modules/%s/%s/%s/", *params.Namespace, *params.Name, *params.Provider)
	dirs, err := svc.storageSvc.ListDirs(ctxId, dirPath)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "ListModuleVersions", params)
	}

	var versions []Version
	for _, dir := range *dirs {
		dir = strings.TrimPrefix(dir, dirPath)
		dir = strings.TrimSuffix(dir, "/")
		versions = append(versions, Version{Version: dir})
	}

	result := ListModuleVersionsResponse{
		Modules: []ModuleVersions{{Versions: versions}},
	}

	svc.logger.Debug(fmt.Sprintf("%s ListModuleVersions() return", ctxId),
		zap.Reflect("result", result),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Name),
		zap.String("provider", *params.Provider),
	)
	return &result, nil
}

func (svc *Modules) GetDownloadUrl(ctxId string, params InputParams) (*string, *ModulesError) {
	svc.logger.Debug(fmt.Sprintf("%s modules.GetDownloadUrl() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Name),
		zap.String("provider", *params.Provider),
		zap.String("version", *params.Version),
	)
	key := fmt.Sprintf("modules/%[1]s/%[2]s/%[3]s/%[4]s/%[1]s-%[2]s-%[3]s-%[4]s.tar.gz", *params.Namespace, *params.Name, *params.Provider, *params.Version)
	result, err := svc.storageSvc.GetDownloadUrl(ctxId, key)

	if err != nil {
		return nil, svc.handleError(ctxId, err, "ListModuleVersions", params)
	}

	svc.logger.Debug(fmt.Sprintf("%s ListModuleVersions() return", ctxId),
		zap.Reflect("result", result),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Name),
		zap.String("provider", *params.Provider),
		zap.String("version", *params.Version),
	)
	return result, nil
}

func (svc *Modules) handleError(ctxId string, err *storage.StorageError, method string, params InputParams, fields ...zap.Field) *ModulesError {
	if err.Code == storage.ErrObjectNotFound {
		fields = append(fields, zap.NamedError("errStorage", err))
		svc.logger.Warn(fmt.Sprintf("%s Modules.%s() Storage.NotFound", ctxId, method),
			fields...,
		)
		return &ModulesError{
			Message: fmt.Sprintf("Error #%d Module %s/%s/%s not found", ErrNotFound, *params.Namespace, *params.Name, *params.Provider),
			Code:    ErrNotFound,
			Err:     err,
		}
	}

	fields = append(fields, zap.NamedError("errStorage", err))
	svc.logger.Warn(fmt.Sprintf("%s Modules.%s() Storage.NotFound", ctxId, method),
		fields...,
	)
	return &ModulesError{
		Message: fmt.Sprintf("Error #%d unknonwn error with module %s/%s/%s", ErrUnknown, *params.Namespace, *params.Name, *params.Provider),
		Code:    ErrUnknown,
		Err:     err,
	}
}
