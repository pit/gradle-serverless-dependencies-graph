package providers

import (
	"regexp"
	"strings"
	"terraform-serverless-private-registry/lib/helpers"
	storage2 "terraform-serverless-private-registry/lib/storage"
	"testing"
)

func TestProvidersNewProviders(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, err := NewProviders(storage, logger)

	if modules == nil {
		t.Error("Providers is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}
}

func TestProvidersListProviderVersions(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, _ := NewProviders(storage, logger)

	providerNamespace := "kvinta"
	providerType := "dbmigration"
	params := ListProviderVersionsInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
	}
	resp, err := modules.ListProviderVersions("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	if resp.Versions == nil {
		t.Error("No modules in response")
	}

	ok := false
	for _, provider := range resp.Versions {
		if provider.Version == "" {
			t.Error("No versions in response.Providers")
		}

		for _, version := range resp.Versions {
			m, err := regexp.MatchString(`\d+\.\d+\.\d+`, version.Version)
			if err != nil {
				t.Error(err)
			}
			if m {
				ok = true
			}
		}
	}
	if !ok {
		t.Error("Wrong resp content")
	}

}

func TestProvidersGetDownload(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	providers, _ := NewProviders(storage, logger)

	providerNamespace := "kvinta"
	providerType := "dbmigration"
	providerVersion := "1.0.6"
	providerOS := "darwin"
	providerArch := "amd64"
	params := GetDownloadInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
		Version:   &providerVersion,
		OS:        &providerOS,
		Arch:      &providerArch,
	}
	resp, err := providers.GetDownload("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	if !strings.HasPrefix(resp.DownloadUrl, "https://terraform-registry-kvinta-io.s3.eu-central-1.amazonaws.com/providers/kvinta/dbmigration/1.0.6/terraform-provider-dbmigration_1.0.6_darwin_amd64.zip?") {
		t.Error("Wrong url format")
	}
}

func TestProvidersGetDownloadUrl404(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	providers, _ := NewProviders(storage, logger)

	providerNamespace := "kvinta"
	providerType := "dbmigration"
	providerVersion := "0.0.0"
	providerOS := "darwin"
	providerArch := "amd64"

	params := GetDownloadInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
		Version:   &providerVersion,
		OS:        &providerOS,
		Arch:      &providerArch,
	}
	resp, err := providers.GetDownload("0000", params)

	if resp != nil {
		t.Error("Response is not nil")
	}
	if err == nil {
		t.Error("Should be error")
	}
}
