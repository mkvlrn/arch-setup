.DEFAULT_GOAL := lint

.PHONY: setup lint format format-check sync-branch

setup:
	mise trust --yes
	mise install
	lefthook install

lint:
	@shellcheck config.sh install.sh verify.sh .github/workflows/scripts/*.sh
	@bash -n config.sh install.sh verify.sh .github/workflows/scripts/*.sh

format:
	@shfmt -w config.sh install.sh verify.sh .github/workflows/scripts/*.sh

format-check:
	@shfmt -d config.sh install.sh verify.sh .github/workflows/scripts/*.sh

sync-branch:
	@if command -v mise >/dev/null 2>&1; then mise prune -y; fi
	@lefthook install

%:
	@:
