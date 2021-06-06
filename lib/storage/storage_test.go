package storage

import (
	"net/http"
	"terraform-serverless-private-registry/lib/helpers"
	"testing"
)

func TestNewStorage(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, err := NewStorage(bucketName, logger)

	if storage == nil {
		t.Error("storageSvc is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}
}

func TestListDirs(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)

	resp, err := storage.ListDirs("0000", "modules/")

	if err != nil {
		t.Error(err)
	}
	if len(*resp) == 0 {
		t.Error("Response is empty")
	}

	ok := false
	for _, item := range *resp {
		if item == "modules/apps/" {
			ok = true
		}
	}
	if !ok {
		t.Error("Wrong resp content")
	}
}

func TestGetDownloadUrl(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)

	key := "modules/apps/cryptopro/k8s/0.1.134/apps-cryptopro-k8s-0.1.134.tar.gz"
	resp, err := storage.GetDownloadUrl("0000", key)

	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("Response is nil")
	}

	httpResp, httpErr := http.Get(*resp)
	if httpErr != nil {
		t.Error(httpErr)
	}
	if httpResp.StatusCode != 200 {
		t.Errorf("Wrong response status %d", httpResp.StatusCode)
	}
}

func TestGetDownloadUrlNotFound(t *testing.T) {
	bucketName := "terraform-registry-kvinta-io"
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)

	key := "not_exist"
	resp, err := storage.GetDownloadUrl("0000", key)

	if err == nil {
		t.Error("Error is nil")
	}
	if resp != nil {
		t.Error("Response is not empty")
	}

	if err.Code != ErrObjectNotFound {
		t.Error("Wrong error code", err)
	}
}
