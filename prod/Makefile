DOCKER_CMD=REPO_NAME=$(REPO_NAME) REPO_ROOT=$(REPO_ROOT) DOMAIN_PROD=$(DOMAIN_PROD) API_DOMAIN_PROD=$(API_DOMAIN_PROD) docker compose -f $(REPO_ROOT)/prod/docker-compose.yml
DOCKER_ROLLOUT_CMD=REPO_NAME=$(REPO_NAME) REPO_ROOT=$(REPO_ROOT) DOMAIN_PROD=$(DOMAIN_PROD) API_DOMAIN_PROD=$(API_DOMAIN_PROD) docker rollout -f $(REPO_ROOT)/prod/docker-compose.yml

docker: ## setup docker on new server
	@sudo sysctl -w net.core.rmem_max=7500000
	@sudo sysctl -w net.core.wmem_max=7500000
	@which docker > /dev/null || wget -qO- https://get.docker.com/ | sh
	@sudo usermod -aG docker $(shell whoami)
	@mkdir -p ~/.docker/cli-plugins && curl https://raw.githubusercontent.com/wowu/docker-rollout/master/docker-rollout -o ~/.docker/cli-plugins/docker-rollout && chmod +x ~/.docker/cli-plugins/docker-rollout

shell-db: ## login running db container
	@docker load < $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-db.tar.gz
	@docker run --rm -it -v $(REPO_ROOT)/db:/db -w /db $(REPO_NAME)-prod-db:latest bash

shell-frontend: ## login running frontend container
	@$(DOCKER_CMD) exec -w /usr/share/nginx/html frontend sh

init: ## init on new server, this also runs `make migrate` to create db first.
	@docker load < $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-backend.tar.gz
	@docker load < $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-frontend.tar.gz
	@docker load < $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-db.tar.gz
	@docker run --rm -it -v $(REPO_ROOT)/db:/db -w /db $(REPO_NAME)-prod-db:latest /db/migrate
	@$(DOCKER_CMD) up -d

update-backend: ## load image and docker-rollout
	@docker load < $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-backend.tar.gz
	@$(DOCKER_ROLLOUT_CMD) backend

update-frontend: ## load image and docker-rollout
	@docker load < $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-frontend.tar.gz
	@$(DOCKER_ROLLOUT_CMD) frontend

prune: ## prune all images and volumes
	@docker system prune --all --volumes --force
