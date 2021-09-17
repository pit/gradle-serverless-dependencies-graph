gobuildcmd = GOSUMDB="off" GOPRIVATE=github.com/pit/terraform-serverless-private-registry CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  go build -installsuffix cgo -ldflags "-X main.version=`cat version` -X main.builddate=`date -u +.%Y%m%d.%H%M%S` -w -s"

.PHONY: dep
dep:
	go mod download

.PHONY: build-lambda
build:
	#$(gobuildcmd) -o bin/authorizer lambda/authorizer/*.go

	# lambda for index response
	$(gobuildcmd) -o bin/default lambda/default/*.go
	#$(gobuildcmd) -o bin/index lambda/index/*.go
	#$(gobuildcmd) -o bin/discovery lambda/discovery/*.go

	$(gobuildcmd) -o bin/repo-batch-insert-put lambda/repo-batch-insert-put/*.go

pack:
	mkdir -p dist
	#zip -j dist/authorizer.zip bin/authorizer

	zip -j dist/default.zip bin/default
	#zip -j dist/index.zip bin/index
	#zip -j dist/discovery.zip bin/discovery
	zip -j dist/repo-batch-insert-put.zip bin/repo-batch-insert-put
