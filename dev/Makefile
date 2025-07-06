# avoid .env file if possible, it's error prone, and sooner or later you'll find yourself unable to run the project at all
DOMAIN_PROD?=example.com
IMAGE_GO?=golang:1.24.4-bookworm
VERSION_STATICCHECK?=v0.6.1
VERSION_REVIVE?=v1.10.0
VERSION_XX?=1.6.1
IMAGE_GO_PROD?=gcr.io/distroless/base-debian12
IMAGE_NODE?=node:24.2.0-bookworm
IMAGE_NGINX?=nginx:1.28.0-alpine
IMAGE_UV?=ghcr.io/astral-sh/uv:python3.13-bookworm-slim
# git@github.com:user/repo.git or https://github.com/user/repo.git
# => github.com/user/repo
REPO?=$(shell git config --get remote.origin.url | sed -E 's|git@||; s|https://||; s|\.git$$||; s|:|/|')
# github.com/user/repo => repo
REPO_NAME?=$(shell echo $(REPO) | sed -E 's|.*/||')
REPO_ROOT?=$(shell git rev-parse --show-toplevel)
# configure in ~/.ssh/config
REMOTE_ROOT?=nerd:/home/yk/repos/$(REPO_NAME)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || true)
LDFLAGS?="-X 'github.com/oseau/web.VersionString=dev' -X 'github.com/oseau/web/cmd/http.URLFrontend=https://$(REPO_NAME).orb.local'"
LDFLAGS_PROD?="-X 'github.com/oseau/web.VersionString=$(GIT_COMMIT)$(GIT_DIRTY)' -X 'github.com/oseau/web/cmd/http.URLFrontend=https://$(DOMAIN_PROD)'"
API_URL?=https://api.$(REPO_NAME).orb.local
API_URL_PROD?=https://api.$(DOMAIN_PROD)
DOCKER_CMD=COMPOSE_BAKE=true REPO_NAME=$(REPO_NAME) REPO_ROOT=$(REPO_ROOT) IMAGE_GO=$(IMAGE_GO) VERSION_STATICCHECK=$(VERSION_STATICCHECK) VERSION_REVIVE=$(VERSION_REVIVE) LDFLAGS=$(LDFLAGS) IMAGE_NODE=$(IMAGE_NODE) API_URL=$(API_URL) IMAGE_UV=$(IMAGE_UV) docker compose -f $(REPO_ROOT)/dev/docker-compose.yml
IMAGE_TQDM?=tqdm/tqdm:4.67.1 # to display progress bar

dev: ## start dev server for both backend & frontend
	@$(DOCKER_CMD) up --build --renew-anon-volumes --remove-orphans

shell-backend: ## login running backend container
	@$(DOCKER_CMD) exec backend bash

shell-db: ## login running db container
	@$(DOCKER_CMD) exec db bash

shell-frontend: ## login running frontend container
	@$(DOCKER_CMD) exec frontend bash

build: ## build for deploy
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/backend/build.Dockerfile --build-arg IMAGE_GO=$(IMAGE_GO) --build-arg IMAGE_GO_PROD=$(IMAGE_GO_PROD) --build-arg VERSION_XX=$(VERSION_XX) --build-arg LDFLAGS=$(LDFLAGS_PROD) --tag $(REPO_NAME)-prod-backend $(REPO_ROOT)
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/frontend/build.Dockerfile --build-arg IMAGE_NODE=$(IMAGE_NODE) --build-arg API_URL=$(API_URL_PROD) --build-arg IMAGE_NGINX=$(IMAGE_NGINX) --tag $(REPO_NAME)-prod-frontend $(REPO_ROOT)
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/db/build.Dockerfile --build-arg IMAGE_UV=$(IMAGE_UV) --tag $(REPO_NAME)-prod-db $(REPO_ROOT)

rsync: ## rsync built images to remote server
	@docker save $(REPO_NAME)-prod-backend:latest | docker run --rm -i $(IMAGE_TQDM) --bytes --total $(shell docker image inspect "$(REPO_NAME)-prod-backend:latest" --format='{{.Size}}') | gzip > $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-backend.tar.gz
	@docker save $(REPO_NAME)-prod-frontend:latest | docker run --rm -i $(IMAGE_TQDM) --bytes --total $(shell docker image inspect "$(REPO_NAME)-prod-frontend:latest" --format='{{.Size}}') | gzip > $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-frontend.tar.gz
	@docker save $(REPO_NAME)-prod-db:latest | docker run --rm -i $(IMAGE_TQDM) --bytes --total $(shell docker image inspect "$(REPO_NAME)-prod-db:latest" --format='{{.Size}}') | gzip > $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-db.tar.gz
	@rsync -azvhP --delete --mkpath --checksum $(REPO_ROOT)/tmp/*.tar.gz $(REMOTE_ROOT)/tmp/
