gobuildcmd = GOSUMDB="off" GOPRIVATE=github.com/pit/terraform-serverless-private-registry CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  go build -installsuffix cgo -ldflags "-X main.version=`cat version` -X main.builddate=`date -u +.%Y%m%d.%H%M%S` -w -s"

.PHONY: dep
dep:
	go mod download

.PHONY: build-lambda
build:
	$(gobuildcmd) -o bin/authorizer lambda/authorizer/*.go

	# lambda for index response
	$(gobuildcmd) -o bin/default lambda/default/*.go
	$(gobuildcmd) -o bin/index lambda/index/*.go

	$(gobuildcmd) -o bin/api-repository-batch-insert lambda/api-repository-batch-insert/*.go
	$(gobuildcmd) -o bin/web-dependencies-list-by-parent lambda/web-dependencies-list-by-parent/*.go
	$(gobuildcmd) -o bin/web-dependencies-list-by-repo lambda/web-dependencies-list-by-repo/*.go
	$(gobuildcmd) -o bin/web-repositories-list-by-parent lambda/web-repositories-list-by-parent/*.go
	$(gobuildcmd) -o bin/web-repositories-list-by-dep lambda/web-repositories-list-by-dep/*.go

pack:
	mkdir -p dist
	zip -j dist/authorizer.zip bin/authorizer

	zip -j dist/default.zip bin/default
	zip -j dist/index.zip bin/index

	zip -j dist/api-repository-batch-insert.zip bin/api-repository-batch-insert
	zip -j dist/web-dependencies-list-by-parent.zip bin/web-dependencies-list-by-parent
	zip -j dist/web-dependencies-list-by-repo.zip bin/web-dependencies-list-by-repo
	zip -j dist/web-repositories-list-by-parent.zip bin/web-repositories-list-by-parent
	zip -j dist/web-repositories-list-by-dep.zip bin/web-repositories-list-by-dep

