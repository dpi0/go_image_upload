BINARY_NAME := "server"
MAIN_PATH := "./cmd/server"
BINARY_PATH := "./bin"

air:
  air

tidy:
  go mod tidy -v
  go fmt ./...

build:
	GOARCH=amd64 GOOS=darwin go build -o {{BINARY_PATH}}/{{BINARY_NAME}}-darwin {{MAIN_PATH}}
	GOARCH=amd64 GOOS=linux go build -o {{BINARY_PATH}}/{{BINARY_NAME}}-linux {{MAIN_PATH}}
	GOARCH=amd64 GOOS=windows go build -o {{BINARY_PATH}}/{{BINARY_NAME}}-windows {{MAIN_PATH}}

run: build
	{{BINARY_PATH}}/{{BINARY_NAME}}

clean:
	go clean
	rm {{BINARY_PATH}}/{{BINARY_NAME}}-darwin
	rm {{BINARY_PATH}}/{{BINARY_NAME}}-linux
	rm {{BINARY_PATH}}/{{BINARY_NAME}}-windows

test:
	go test ./...