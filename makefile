.PHONY: default run run-with-docs build test doc clean update-docs hard-clean
# Variables
APP_NAME = "app"

# Determine the platform
ifdef ComSpec
    RM = if exist "$(APP_NAME)$(APP_EXT)" del /Q /F
    RMDIR = if exist
    APP_EXT = .exe
else
    RM = rm -rf
    RMDIR = rm -rf
    APP_EXT =
endif

# Tasks
default: run

run:
	@go run ./cmd/main.go

run-with-docs: docs run

build: update-docs
	@go build -o $(APP_NAME) ./cmd/main.go

test:
	@go test ./...

docs:
	@swag init -g cmd/main.go

clean:
	@$(RM) "$(APP_NAME)$(APP_EXT)"
	@if exist docs rmdir /S /Q docs

update-docs: clean docs

hard-clean: clean
	@if exist db_data rmdir /S /Q db_data
