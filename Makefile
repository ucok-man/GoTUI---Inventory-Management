# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^//'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# COMMAND
# ==================================================================================== #

## app/run: run the application without creating binary
.PHONY: app/run
app/run:
	@go run ./src/cmd/inventory-management/main.go

## app/build: build the binary version of this application
.PHONY: app/build
app/build:
	@go build -o ./build/app ./src/cmd/inventory-management/main.go

## app/start: run the output binary
.PHONY: app/start
app/start:
	@./build/app