SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

# avoid .env file if possible, it's error prone, and sooner or later you'll find yourself unable to run the project at all
DOMAIN_PROD?=example.com
API_DOMAIN_PROD?=api.example.com
# same debian version as prod distroless image
IMAGE_GO?=golang:1.25.0-bookworm
VERSION_STATICCHECK?=v0.6.1
VERSION_REVIVE?=v1.10.0
VERSION_XX?=1.6.1
IMAGE_GO_PROD?=gcr.io/distroless/base-debian12
IMAGE_NODE?=node:24.5.0-bookworm
IMAGE_NGINX?=nginx:1.29.0-alpine
IMAGE_UV?=ghcr.io/astral-sh/uv:python3.13-bookworm-slim
IMAGE_REDIS?=redis:8.2.0-alpine
# git@github.com:user/repo.git or https://github.com/user/repo.git
# => github.com/user/repo
REPO?=$(shell git config --get remote.origin.url | sed -E 's|git@||; s|https://||; s|\.git$$||; s|:|/|')
# github.com/user/repo => repo
REPO_NAME?=$(shell echo $(REPO) | sed -E 's|.*/||')
# REPO_ROOT is the absolute path to the root of the repo, it's dynamic, different between dev and prod
REPO_ROOT?=$(shell git rev-parse --show-toplevel)
# configure in ~/.ssh/config
REMOTE_ROOT?=nerd:~/repos/$(REPO_NAME)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
LDFLAGS?="-X 'github.com/oseau/web.VersionString=dev' -X 'github.com/oseau/web/http.URLFrontend=https://$(REPO_NAME).orb.local'"
LDFLAGS_PROD?="-X 'github.com/oseau/web.VersionString=$(GIT_COMMIT)$(GIT_DIRTY)' -X 'github.com/oseau/web/http.URLFrontend=https://$(DOMAIN_PROD)'"
API_URL?=https://api.$(REPO_NAME).orb.local
API_URL_PROD?=https://$(API_DOMAIN_PROD)
DOCKER_CMD=COMPOSE_BAKE=true REPO_NAME=$(REPO_NAME) REPO_ROOT=$(REPO_ROOT) IMAGE_GO=$(IMAGE_GO) VERSION_STATICCHECK=$(VERSION_STATICCHECK) VERSION_REVIVE=$(VERSION_REVIVE) LDFLAGS=$(LDFLAGS) IMAGE_NODE=$(IMAGE_NODE) API_URL=$(API_URL) IMAGE_UV=$(IMAGE_UV) IMAGE_REDIS=$(IMAGE_REDIS) docker compose -f $(REPO_ROOT)/dev/docker-compose.yml
IMAGE_TQDM?=tqdm/tqdm:4.67.1 # to display progress bar


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
