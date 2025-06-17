# Read .env variables
# include .env
# .EXPORT_ALL_VARIABLES:

.PHONY: ${wildcard *}

default:
	@echo "Available targets:"
	@grep '^[^ ]*:.*#' Makefile | awk -F'[:#]' '{printf "  %-20s %s\n", $$1, $$3}'
.PHONY: default


run: # Run the game
	love .

build: # Build the game
	mkdir -p bin
	zip -9 -r bin/game.love main.lua assets/
	mkdir -p public
	npx love.js -t luagame bin/game.love public/
