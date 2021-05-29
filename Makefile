gobuildcmd = GOSUMDB="off" GOPRIVATE=github.com/pit/terraform-serverless-private-registry CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  go build -installsuffix cgo -ldflags "-X main.version=`cat version` -X main.builddate=`date -u +.%Y%m%d.%H%M%S` -w -s"

.PHONY: dep
dep:
	go mod download

.PHONY: build-lambda
build:
    # lambda for index response
	$(gobuildcmd) -o bin/default lambda/default/main.go
	$(gobuildcmd) -o bin/index lambda/index/main.go
	$(gobuildcmd) -o bin/discovery lambda/discovery/main.go

	# https://www.terraform.io/docs/internals/module-registry-protocol.html
	$(gobuildcmd) -o bin/modules-list lambda/modules-list/main.go
	$(gobuildcmd) -o bin/modules-search lambda/modules-search/main.go
	$(gobuildcmd) -o bin/modules-versions lambda/modules-versions/main.go
	$(gobuildcmd) -o bin/modules-download lambda/modules-download/main.go
	$(gobuildcmd) -o bin/modules-latest-version lambda/modules-latest-version/main.go
	$(gobuildcmd) -o bin/modules-get lambda/modules-get/main.go

	# https://www.terraform.io/docs/internals/provider-registry-protocol.html
	$(gobuildcmd) -o bin/providers-versions lambda/providers-versions/main.go
	$(gobuildcmd) -o bin/providers-download lambda/providers-download/main.go
