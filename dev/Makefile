DOMAIN_PROD?=example.com
IMAGE_GO?=golang:1.24.4-bookworm
VERSION_STATICCHECK?=v0.6.1
VERSION_REVIVE?=v1.10.0
VERSION_XX?=1.6.1
IMAGE_GO_PROD?=gcr.io/distroless/base-debian12
IMAGE_NODE?=node:24.2.0-bookworm
IMAGE_NGINX?=nginx:1.28.0-alpine
# git@github.com:user/repo.git or https://github.com/user/repo.git
# => github.com/user/repo
REPO?=$(shell git config --get remote.origin.url | sed -E 's|git@||; s|https://||; s|\.git$$||; s|:|/|')
# github.com/user/repo => repo
REPO_NAME?=$(shell echo $(REPO) | sed -E 's|.*/||')
REPO_ROOT?=$(shell git rev-parse --show-toplevel)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
LDFLAGS?="-X 'github.com/oseau/web.VersionString=dev' -X 'github.com/oseau/web/cmd/http.URLFrontend=https://$(REPO_NAME).orb.local'"
LDFLAGS_PROD?="-X 'github.com/oseau/web.VersionString=$(GIT_COMMIT)$(GIT_DIRTY)' -X 'github.com/oseau/web/cmd/http.URLFrontend=https://$(DOMAIN_PROD)'"
DOCKER_CMD=COMPOSE_BAKE=true REPO_NAME=$(REPO_NAME) REPO_ROOT=$(REPO_ROOT) IMAGE_GO=$(IMAGE_GO) VERSION_STATICCHECK=$(VERSION_STATICCHECK) VERSION_REVIVE=$(VERSION_REVIVE) LDFLAGS=$(LDFLAGS) IMAGE_NODE=$(IMAGE_NODE) docker compose -f $(REPO_ROOT)/dev/docker-compose.yml

dev: ## start dev server for both backend & frontend
	@$(DOCKER_CMD) up --build --renew-anon-volumes --remove-orphans

shell-backend: ## login running backend server
	@$(DOCKER_CMD) exec backend bash

shell-frontend: ## login running frontend server
	@$(DOCKER_CMD) exec frontend bash

build: ## build for deploy
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/backend/build.Dockerfile --build-arg IMAGE_GO=$(IMAGE_GO) --build-arg IMAGE_GO_PROD=$(IMAGE_GO_PROD) --build-arg VERSION_XX=$(VERSION_XX) --build-arg LDFLAGS=$(LDFLAGS_PROD) --tag backend-$(REPO_NAME)-prod $(REPO_ROOT)
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/frontend/build.Dockerfile --build-arg IMAGE_NODE=$(IMAGE_NODE) --build-arg IMAGE_NGINX=$(IMAGE_NGINX) --tag frontend-$(REPO_NAME)-prod $(REPO_ROOT)
