# Show available targets
.PHONY: default
default:
	@echo Available targets:
	@grep "^[^ ]*:" Makefile | grep "#" | awk -F'[:#]' '{printf "  %-20s %s\n", $$1, $$3}'

.PHONY: run
run: # run the program
	nixGLIntel go run .

.PHONY: build
build: # build wasm
	env GOOS=js GOARCH=wasm go build -o public/game.wasm
