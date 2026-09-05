.DEFAULT_GOAL := lint

BASH_SCRIPTS := main.sh $(wildcard misc/*.sh steps/*.sh)
POSIX_SCRIPTS := $(wildcard .github/workflows/scripts/*.sh secrets/*.sh) index.html
SCRIPTS := $(BASH_SCRIPTS) $(POSIX_SCRIPTS)

.PHONY: setup dev lint format test

setup:
	mise trust --yes
	mise install
	lefthook install

dev:
	@bash main.sh $(ARGS)

lint:
	@set -e; for script in $(BASH_SCRIPTS); do bash -n "$$script"; done
	@set -e; for script in $(POSIX_SCRIPTS); do sh -n "$$script"; done
	@shellcheck -x $(SCRIPTS)
	@shfmt -d $(SCRIPTS)

format:
	@shfmt -w $(SCRIPTS)

test:
	@bash misc/test.sh
	@bash misc/test-launcher.sh
