package providers

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"terraform-serverless-private-registry/lib/storage"
)

type Providers struct {
	storageSvc *storage.Storage
	logger     *zap.Logger
}

type ListProviderVersionsOutput struct {
	Versions []ProviderVersion `json:"versions"`
}

type ProviderVersion struct {
	Version   string             `json:"version"`
	Protocols []string           `json:"protocols"`
	Platforms []ProviderPlatform `json:"platforms"`
}
type ProviderPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type ListProviderVersionsInput struct {
	Namespace *string
	Type      *string
}

type GetDownloadInput struct {
	Namespace *string
	Type      *string
	Version   *string
	OS        *string
	Arch      *string
}

type GpgPublicKey struct {
	KeyId          string `json:"key_id"`
	AsciiArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
	Source         string `json:"source"`
	SourceUrl      string `json:"source_url"`
}

type SigningKeys struct {
	GpgPublicKeys []GpgPublicKey `json:"gpg_public_keys"`
}

type GetDownloadOutput struct {
	Protocols     []string    `json:"protocols"`
	OS            string      `json:"os"`
	Arch          string      `json:"arch"`
	Filename      string      `json:"filename"`
	DownloadUrl   string      `json:"download_url"`
	ShaSumsUrl    string      `json:"shasums_url"`
	ShaSumsSigUrl string      `json:"shasums_signature_url"`
	ShaSum        string      `json:"shasum"`
	SigningKeys   SigningKeys `json:"signing_keys"`
}

const (
	ErrNotFound = iota
	ErrUnknown
)

type ProvidersError struct {
	Message string
	Code    int
	Err     error
}

func (s ProvidersError) Error() string {
	panic("implement me")
}

func NewProviders(storage *storage.Storage, log *zap.Logger) (*Providers, error) {
	return &Providers{
		storageSvc: storage,
		logger:     log,
	}, nil
}

func (svc *Providers) ListProviderVersions(ctxId string, params ListProviderVersionsInput) (*ListProviderVersionsOutput, *ProvidersError) {
	svc.logger.Debug("providers.ListProviderVersionsOutput() called",
		zap.String("ctxId", ctxId),
		zap.Reflect("params", params),
	)
	dirPath := fmt.Sprintf("providers/%s/%s/", *params.Namespace, *params.Type)
	dirs, err := svc.storageSvc.ListDirs(ctxId, dirPath)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "ListProviderVersionsOutput() error", *params.Namespace, *params.Type)
	}

	var versions []ProviderVersion
	for _, dir := range *dirs {
		dir = strings.TrimPrefix(dir, dirPath)
		dir = strings.TrimSuffix(dir, "/")
		version := ProviderVersion{
			Version:   dir,
			Protocols: []string{"4.0"},
			Platforms: []ProviderPlatform{
				ProviderPlatform{
					OS:   "",
					Arch: "",
				},
			},
		}
		versions = append(versions, version)
	}

	result := ListProviderVersionsOutput{
		Versions: versions,
	}

	svc.logger.Debug(fmt.Sprintf("%s ListProviderVersions() return", ctxId),
		zap.Reflect("result", result),
		zap.String("namespace", *params.Namespace),
		zap.String("type", *params.Type),
	)
	return &result, nil
}

func (svc *Providers) GetDownload(ctxId string, params GetDownloadInput) (*GetDownloadOutput, *ProvidersError) {
	svc.logger.Debug(fmt.Sprintf("%s modules.GetDownloadUrl() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("type", *params.Type),
		zap.String("version", *params.Version),
		zap.String("os", *params.OS),
		zap.String("arch", *params.Arch),
	)

	downloadFileName := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_%[3]s_%[4]s.zip", *params.Type, *params.Version, *params.OS, *params.Arch)
	downloadUrlKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_%[4]s_%[5]s.zip", *params.Namespace, *params.Type, *params.Version, *params.OS, *params.Arch)
	downloadUrl, err := svc.storageSvc.GetDownloadUrl(ctxId, downloadUrlKey, downloadFileName)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "GetDownload", *params.Namespace, *params.Type)
	}
	downloadMetadata, err := svc.storageSvc.GetMetadata(ctxId, downloadUrlKey)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "GetDownload", *params.Namespace, *params.Type)
	} else {
		svc.logger.Debug("GetDownload.Metadata",
			zap.Reflect("metadata", downloadMetadata),
		)
	}

	shaSumsFileName := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_SHA256SUMS", *params.Type, *params.Version)
	shaSumsKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS", *params.Namespace, *params.Type, *params.Version)
	shaSumsUrl, err := svc.storageSvc.GetDownloadUrl(ctxId, shaSumsKey, shaSumsFileName)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "GetDownload", *params.Namespace, *params.Type)
	}

	shaSumsSigFileName := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_SHA256SUMS.sig", *params.Type, *params.Version)
	shaSumsSigKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS.sig", *params.Namespace, *params.Type, *params.Version)
	shaSumsSigUrl, err := svc.storageSvc.GetDownloadUrl(ctxId, shaSumsSigKey, shaSumsSigFileName)
	if err != nil {
		return nil, svc.handleError(ctxId, err, "GetDownload", *params.Namespace, *params.Type)
	}

	shaSum := (*downloadMetadata)["sha256"]

	result := GetDownloadOutput{
		Protocols:     []string{"4.0"},
		OS:            *params.OS,
		Arch:          *params.Arch,
		Filename:      downloadFileName,
		DownloadUrl:   *downloadUrl,
		ShaSumsUrl:    *shaSumsUrl,
		ShaSumsSigUrl: *shaSumsSigUrl,
		ShaSum:        shaSum,
		SigningKeys: SigningKeys{
			GpgPublicKeys: []GpgPublicKey{
				{
					KeyId:          "",
					AsciiArmor:     "",
					TrustSignature: "",
					Source:         "",
					SourceUrl:      "",
				},
			},
		},
	}

	svc.logger.Debug(fmt.Sprintf("%s ListModuleVersions() return", ctxId),
		zap.Reflect("result", result),
		zap.String("namespace", *params.Namespace),
		zap.String("type", *params.Type),
		zap.String("version", *params.Version),
		zap.String("os", *params.OS),
		zap.String("arch", *params.Arch),
	)
	return &result, nil
}

func (svc *Providers) handleError(ctxId string, err *storage.StorageError, method string, providerNamespace string, providerType string, fields ...zap.Field) *ProvidersError {
	if err.Code == storage.ErrObjectNotFound {
		fields = append(fields, zap.NamedError("errStorage", err))
		svc.logger.Warn(fmt.Sprintf("%s Modules.%s() Storage.NotFound", ctxId, method),
			fields...,
		)
		return &ProvidersError{
			Message: fmt.Sprintf("Error #%d Provider %s/%s not found", ErrNotFound, providerNamespace, providerType),
			Code:    ErrNotFound,
			Err:     err,
		}
	}

	fields = append(fields, zap.NamedError("errStorage", err))
	svc.logger.Warn(fmt.Sprintf("%s Modules.%s() Storage.NotFound", ctxId, method),
		fields...,
	)
	return &ProvidersError{
		Message: fmt.Sprintf("Error #%d unknonwn error with provider %s/%s", ErrUnknown, providerNamespace, providerType),
		Code:    ErrUnknown,
		Err:     err,
	}
}
