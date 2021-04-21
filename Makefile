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

