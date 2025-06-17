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
