package modules

import (
	"regexp"
	"strings"
	"terraform-serverless-private-registry/lib/helpers"
	storage2 "terraform-serverless-private-registry/lib/storage"
	"testing"
)

func TestModulesNewModules(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, err := NewModules(storage, logger)

	if modules == nil {
		t.Error("Modules is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}
}

func TestModulesListModuleVersions(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, _ := NewModules(storage, logger)

	namespace := "apps"
	name := "cryptopro"
	provider := "k8s"
	params := InputParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
	}
	resp, err := modules.ListModuleVersions("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	if resp.Modules == nil {
		t.Error("No modules in response")
	}

	ok := false
	for _, module := range resp.Modules {
		if module.Versions == nil {
			t.Error("No versions in response.Modules")
		}

		for _, version := range module.Versions {
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

func TestModulesGetDownloadUrl(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, _ := NewModules(storage, logger)

	namespace := "apps"
	name := "cryptopro"
	provider := "k8s"
	version := "0.1.134"
	params := InputParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
		Version:   &version,
	}
	resp, err := modules.GetDownloadUrl("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	if !strings.HasPrefix(*resp, "https://terraform-registry-kvinta-io.s3.eu-central-1.amazonaws.com/modules/apps/cryptopro/k8s/0.1.134/apps-cryptopro-k8s-0.1.134.tar.gz?") {
		t.Error("Wrong url format")
	}
}

func TestModulesGetDownloadUrl404(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, _ := NewModules(storage, logger)

	namespace := "apps"
	name := "cryptopro"
	provider := "k8s"
	version := "0.0.0"
	params := InputParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
		Version:   &version,
	}
	resp, err := modules.GetDownloadUrl("0000", params)

	if resp != nil {
		t.Error("Response is not nil")
	}
	if err == nil {
		t.Error("Should be error")
	}
}
