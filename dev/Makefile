dev: ## start dev server for both backend & frontend (db will be created/migrated)
	@$(DOCKER_CMD) up --build --renew-anon-volumes --remove-orphans

lint: ## lint the code
	@$(DOCKER_CMD) exec backend go mod tidy -v
	@$(DOCKER_CMD) exec backend go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test=false ./... # skip test files

migrate: ## migrate/create db
	@docker run --rm -it -v $(REPO_ROOT):/usr/src/app -w /usr/src/app/db $(REPO_NAME)-db /usr/src/app/db/migrate

shell-backend: ## login running backend container
	@$(DOCKER_CMD) exec backend bash

shell-db: ## start a new db container (uses compose-built image)
	@docker run --rm -it -v $(REPO_ROOT):/usr/src/app -v $(REPO_ROOT)/data:/data -w /usr/src/app/db $(REPO_NAME)-db bash

shell-frontend: ## login running frontend container
	@$(DOCKER_CMD) exec -w /usr/src/app/frontend frontend bash

shell-redis: ## login running redis container
	@$(DOCKER_CMD) exec redis redis-cli

build: ## build for deploy
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/backend/build.Dockerfile --build-arg IMAGE_GO=$(IMAGE_GO) --build-arg IMAGE_GO_PROD=$(IMAGE_GO_PROD) --build-arg VERSION_XX=$(VERSION_XX) --build-arg LDFLAGS=$(LDFLAGS_PROD) --tag $(REPO_NAME)-prod-backend $(REPO_ROOT)
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/frontend/build.Dockerfile --build-arg IMAGE_NODE=$(IMAGE_NODE) --build-arg API_URL=$(API_URL_PROD) --build-arg IMAGE_NGINX=$(IMAGE_NGINX) --tag $(REPO_NAME)-prod-frontend $(REPO_ROOT)
	@docker buildx build --platform linux/amd64 --file $(REPO_ROOT)/dev/db/build.Dockerfile --build-arg IMAGE_UV=$(IMAGE_UV) --tag $(REPO_NAME)-prod-db $(REPO_ROOT)

rsync: ## rsync built images to remote server
	@docker save $(REPO_NAME)-prod-backend:latest | docker run --rm -i $(IMAGE_TQDM) --bytes --total $(shell docker image inspect "$(REPO_NAME)-prod-backend:latest" --format='{{.Size}}') | gzip > $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-backend.tar.gz
	@docker save $(REPO_NAME)-prod-frontend:latest | docker run --rm -i $(IMAGE_TQDM) --bytes --total $(shell docker image inspect "$(REPO_NAME)-prod-frontend:latest" --format='{{.Size}}') | gzip > $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-frontend.tar.gz
	@docker save $(REPO_NAME)-prod-db:latest | docker run --rm -i $(IMAGE_TQDM) --bytes --total $(shell docker image inspect "$(REPO_NAME)-prod-db:latest" --format='{{.Size}}') | gzip > $(REPO_ROOT)/tmp/$(REPO_NAME)-prod-db.tar.gz
	@rsync -azvhP --delete --mkpath --checksum --exclude='*.db' --exclude='*.db-shm' --exclude='*.db-wal' --exclude='*.db-journal' --exclude='node_modules' --exclude='secrets' --exclude='.env' --exclude='.venv' --exclude='*.pem' --exclude='*.key' $(REPO_ROOT)/ $(REMOTE_ROOT)

logs: ## tail logs
	@$(DOCKER_CMD) logs -f
