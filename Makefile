.DEFAULT_GOAL := lint

.PHONY: lint format format-check sync-branch

lint:
	@shellcheck config.sh install.sh verify.sh .github/workflows/scripts/*.sh secrets/*.sh
	@bash -n config.sh install.sh verify.sh .github/workflows/scripts/*.sh secrets/*.sh

format:
	@shfmt -w config.sh install.sh verify.sh .github/workflows/scripts/*.sh secrets/*.sh

format-check:
	@shfmt -d config.sh install.sh verify.sh .github/workflows/scripts/*.sh secrets/*.sh

sync-branch:
	@if command -v mise >/dev/null 2>&1; then mise prune -y; fi
	@lefthook install

%:
	@:
