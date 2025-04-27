# Show available targets
.PHONY: default
default:
	@echo Available targets:
	@grep "^[^ ]*:" Makefile | grep "#" | awk -F'[:#]' '{printf "  %-20s %s\n", $$1, $$3}'

.PHONY: run
run: # run the program
	nixGLIntel go run .


.PHONY: build
build: # build to wasm
	# env GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o public/game.wasm
	env GOOS=js GOARCH=wasm go build -o public/game.wasm

wasm-opt:
	wasm-opt -O4 --enable-bulk-memory -o ./public/game.wasm ./public/game.wasm 

tinygo:
	tinygo build -target wasm -o ./public/game.wasm

### tinygo wasm_exec.js
# cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js ./public/

