IMAGE_GO?=golang:1.24.4-bookworm
VERSION_STATICCHECK?=v0.6.1
VERSION_REVIVE?=v1.10.0
# git@github.com:user/repo.git or https://github.com/user/repo.git
# => github.com/user/repo
REPO?=$(shell git config --get remote.origin.url | sed -E 's|git@||; s|https://||; s|\.git$$||; s|:|/|')
# github.com/user/repo => repo
REPO_NAME?=$(shell echo $(REPO) | sed -E 's|.*/||')
PROJECT_ROOT?=$(shell git rev-parse --show-toplevel)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
VERSION_FLAG=-X $(REPO).VersionString=$(GIT_COMMIT)$(GIT_DIRTY)
DOCKER_CMD=COMPOSE_BAKE=true NAME=$(REPO_NAME) ROOT=$(PROJECT_ROOT) IMAGE_GO=$(IMAGE_GO) VERSION_STATICCHECK=$(VERSION_STATICCHECK) VERSION_REVIVE=$(VERSION_REVIVE) docker compose -f $(PROJECT_ROOT)/dev/docker-compose.yml

dev: ## start dev server
	@$(DOCKER_CMD) up --build --renew-anon-volumes

shell: ## login running dev server
	@$(DOCKER_CMD) exec $(REPO_NAME) bash
