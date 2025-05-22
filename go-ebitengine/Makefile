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
	env GOOS=js GOARCH=wasm go build -o public/game.wasm

## worse performance somehow?
# env GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o public/game.wasm
# `-s` Omit the symbol table
# `-w  Omit DWARF debugging information

wasm-opt: # optimize wasm
	wasm-opt -O4 --enable-bulk-memory -o ./public/game.wasm ./public/game.wasm 
	
### size
# wasm-opt -Oz --enable-bulk-memory --dce --vacuum --strip-debug --strip-dwarf --coalesce-locals -o public/game.wasm public/game.wasm
### performance
# wasm-opt -O4 --enable-bulk-memory --vacuum --fast-math --flatten --rereloop --coalesce-locals --enable-mutable-globals -o public/game.wasm public/game.wasm
### simple
# wasm-opt -Oz --enable-bulk-memory -o ./public/game.wasm ./public/game.wasm 
# wasm-opt -Oz --enable-bulk-memory --dce --vacuum --strip-debug --strip-dwarf --fast-math -o ./public/game.wasm ./public/game.wasm 
# wasm-opt -O4 --enable-bulk-memory --flatten --local-cse --dce --vacuum --strip-debug --strip-dwarf --fast-math -o ./public/game.wasm ./public/game.wasm 
# wasm-opt -Oz --enable-bulk-memory -o ./public/game.wasm ./public/game.wasm 

# wasm-opt -O4 --enable-bulk-memory --vacuum --fast-math --flatten --rereloop --coalesce-locals -o ./public/game.wasm ./public/game.wasm 

# tinygo:
# 	tinygo build -target wasm -o ./public/game.wasm

### tinygo wasm_exec.js:
# cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js ./public/

