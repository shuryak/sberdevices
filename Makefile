ifneq (,$(wildcard ./config/.env))
	include ./config/.env
	export
endif

REMOTE_CONNECTION = $(REMOTE_USERNAME)@$(REMOTE_HOSTNAME)

.PHONY: remote-deploy
remote-deploy: image-tar-backend-prod image-tar-web-auth
	@echo ">> Remote deploy to $(REMOTE_CONNECTION)"
	@ssh $(REMOTE_CONNECTION) "docker rmi -f web-auth backend-debug backend-prod && rm -rf ~/.tmp && mkdir -p ~/.tmp"
	@rsync -av --progress ./build ./.tmp ./scripts ./config ./migrations Makefile $(REMOTE_CONNECTION):~/
	@ssh $(REMOTE_CONNECTION) "make load-tars nginx-init web-auth infra prod"

.PHONY: check-scripts-x-flag
check-scripts-x-flag:
	@test -z "$$(find ./scripts -type f ! -perm -u=x)" || sudo chmod +x ./scripts/*

.PHONY: check-gomplate
check-gomplate: check-scripts-x-flag
	@$(CURDIR)/scripts/check-gomplate.sh

.PHONY: nginx-init
nginx-init: check-gomplate
	@$(CURDIR)/scripts/nginx-init.sh

.PHONY: infra
infra: check-scripts-x-flag
	@$(CURDIR)/scripts/infra.sh

.PHONY: web-auth
web-auth: check-scripts-x-flag
	@$(CURDIR)/scripts/web-auth.sh

.PHONY: debug
debug: check-scripts-x-flag
	@$(CURDIR)/scripts/debug.sh

.PHONY: prod
prod: check-scripts-x-flag
	@$(CURDIR)/scripts/prod.sh

.PHONY: certs-renew-dry
certs-renew-dry: check-scripts-x-flag
	@$(CURDIR)/scripts/certs-renew-dry.sh

.PHONY: certs-renew-force
certs-renew-force: check-scripts-x-flag
	@$(CURDIR)/scripts/certs-renew-force.sh

.PHONY: load-tars
load-tars: check-scripts-x-flag
	@$(CURDIR)/scripts/load-tars.sh

.PHONY: image-tar-web-auth
image-tar-web-auth: check-scripts-x-flag
	@OS=linux \
	ARCH=amd64 \
	TAG=web-auth \
	DOCKERFILE=./build/web-auth/Dockerfile \
	TAR_NAME=web-auth-image.tar \
	$(CURDIR)/scripts/image-tar.sh

.PHONY: image-tar-backend-debug
image-tar-backend-debug: check-scripts-x-flag
	@OS=linux \
	ARCH=amd64 \
	TAG=backend-debug \
	DOCKERFILE=./build/debug/Dockerfile \
	TAR_NAME=backend-debug-image.tar \
	$(CURDIR)/scripts/image-tar.sh

.PHONY: image-tar-backend-prod
image-tar-backend-prod: check-scripts-x-flag
	@OS=linux \
	ARCH=amd64 \
	TAG=backend-prod \
	DOCKERFILE=./build/prod/Dockerfile \
	TAR_NAME=backend-prod-image.tar \
	$(CURDIR)/scripts/image-tar.sh

