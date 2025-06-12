IMAGE_GO?=golang:1.24.4-bookworm
VERSION_STATICCHECK?=v0.6.1
VERSION_REVIVE?=v1.10.0
IMAGE_NODE?=node:24.2.0-bookworm
# git@github.com:user/repo.git or https://github.com/user/repo.git
# => github.com/user/repo
REPO?=$(shell git config --get remote.origin.url | sed -E 's|git@||; s|https://||; s|\.git$$||; s|:|/|')
# github.com/user/repo => repo
REPO_NAME?=$(shell echo $(REPO) | sed -E 's|.*/||')
REPO_ROOT?=$(shell git rev-parse --show-toplevel)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
VERSION_FLAG=-X $(REPO).VersionString=$(GIT_COMMIT)$(GIT_DIRTY)
DOCKER_CMD=COMPOSE_BAKE=true REPO_NAME=$(REPO_NAME) REPO_ROOT=$(REPO_ROOT) IMAGE_GO=$(IMAGE_GO) VERSION_STATICCHECK=$(VERSION_STATICCHECK) VERSION_REVIVE=$(VERSION_REVIVE) IMAGE_NODE=$(IMAGE_NODE) docker compose -f $(REPO_ROOT)/dev/docker-compose.yml

dev: ## start dev server for both backend & frontend
	@$(DOCKER_CMD) up --build --renew-anon-volumes

shell-backend: ## login running backend server
	@$(DOCKER_CMD) exec backend bash

shell-frontend: ## login running frontend server
	@$(DOCKER_CMD) exec frontend bash
