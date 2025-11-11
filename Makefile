.PHONY: make
make:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux \
		go build \
		-ldflags="-s -w -X main.version=@`git describe --tags`" \
		-o ./bin/umami-importer ./main.go

.PHONY: lint
lint:
	deadcode ./...
	modernize ./...
	goimports-reviser -format ./...
	golangci-lint run

.PHONY: updatepackages
updatepackages:
	go mod tidy
	go get ${shell go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all}