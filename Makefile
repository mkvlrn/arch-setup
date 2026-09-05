.DEFAULT_GOAL := build

.PHONY: setup dev lint format format-check test build

setup:
	mise trust --yes
	mise install
	lefthook install

dev:
	@go run . $(ARGS)

lint:
	@golangci-lint run ./...

format:
	@gofumpt -w .

test:
	@go test ./... -cover

# Stamp the actual checkout, including Actions' synthetic PR merge commits.
# Refuse dirty sources: their contents cannot be reproduced from a commit SHA.
build:
	@set -eu; \
		revision=$$(git rev-parse --verify HEAD); \
		state=$$(git --no-optional-locks status --porcelain --untracked-files=normal); \
		if [ -n "$$state" ]; then \
			printf 'Cannot build a revision-aligned installer from a dirty checkout.\n' >&2; \
			exit 1; \
		fi; \
		mkdir -p ./bin; \
		go build -ldflags "-X github.com/mkvlrn/arch-setup/internal/revision.Commit=$$revision" -o ./bin/arch-setup .

%:
	@:
