# Show available targets
.PHONY: default
default:
	@echo Available targets:
	@grep "^[^ ]*:" Makefile | grep "#" | awk -F'[:#]' '{printf "  %-20s %s\n", $$1, $$3}'

.PHONY: run
run: # run the program
	nixVulkanIntel cargo run --features dev

.PHONY: build
build: # build to wasm
	wasm-pack build --target web

.PHONY: release
release: # Build and publish a GitHub release with WASM
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is not set. Usage: make release VERSION=0.1.0"; \
		exit 1; \
	fi
	wasm-pack build --release --target web
	git tag v$(VERSION)
	git push origin v$(VERSION)
	gh release create v$(VERSION) pkg/* --title "v$(VERSION)" --notes "Release WASM build"
