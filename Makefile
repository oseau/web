SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

ifeq ($(shell uname), Darwin)
	include dev/Makefile
else
	include prod/Makefile
endif

# https://www.gnu.org/software/make/manual/html_node/Options-Summary.html
MAKEFLAGS += --always-make

.DEFAULT_GOAL := help
# Modified from http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -Eh '^[a-zA-Z_-]+:.*?##? .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?##? "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
