short-test:
	@go test -v -race -short ./...

test:
	@go test -v -race ./...

format:
	gofumpt -w ./..


install-go-tools:
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %


MOCKGEN_DESTINATION := pkg/mocks
MOCK_SOURCES := pkg/gopherdir/file_service.go pkg/gopherdir/ui.go

mocks: ${MOCK_SOURCES}
	@rm -rf ${MOCKGEN_DESTINATION}/*
	@echo "Generating mocks..."
	@for file in $^; do mockgen -source=$$file -destination=${MOCKGEN_DESTINATION}/$$file; done


BUILD_OUTPUT_DIR?=bin

build-local:
	go build -o $(BUILD_OUTPUT_DIR)/gopherdir-http-local ./cmd/gopherdir-http/...

build-linux-arm-64:
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_OUTPUT_DIR)/gopherdir-http-linux-arm64 ./cmd/gopherdir-http/...
