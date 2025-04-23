# Show available targets
.PHONY: default
default:
	@echo Available targets:
	@grep "^[^ ]*:" Makefile | grep "#" | awk -F'[:#]' '{printf "  %-20s %s\n", $$1, $$3}'

.PHONY: run
run: # run the program
	nixVulkanIntel cargo run

.PHONY: build
build: # build to wasm
	wasm-pack build --target web
