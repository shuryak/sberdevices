ifneq (,$(wildcard .env))
    include .env
endif

CONNECTION = $(REMOTE_USERNAME)@$(REMOTE_HOSTNAME)

# EXAMPLES:
# DEBUG_BIN_NAME = "sberhack_DEBUG"
# DEBUG_DIR = "~/debug"
# DEBUG_WEB_DIR = "~/debug/web"
# DEBUG_SYSTEMCTL_UNIT = "sberdevices-debug.service"
# DEBUG_API_BASE_URL = "https://sdprovider.ru/dev"

# EXAMPLES:
# PROD_BIN_NAME = "sberhack"
# PROD_DIR = "~/prod"
# PROD_WEB_DIR = "~/prod/web"
# PROD_SYSTEMCTL_UNIT = "sberdevices-prod.service"
# PROD_API_BASE_URL = "https://sdprovider.ru"

# REGION: BUILD
.PHONY: build-linux-debug
build-linux-debug:
	@echo ">> Building $(DEBUG_BIN_NAME) (debug) to $(CONNECTION)"
	@env GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o $(DEBUG_BIN_NAME) -v cmd/sberhack/*.go
	@echo "   Done."

.PHONY: build-linux-prod
build-linux-prod:
	@echo ">> Building $(PROD_BIN_NAME) (production) to $(CONNECTION)"
	@env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(PROD_BIN_NAME) -v cmd/sberhack/*.go
	@echo "   Done."
# END REGION: BUILD

# REGION: UPLOADING
.PHONY: upload-debug
upload-debug:
	@echo ">> Uploading $(DEBUG_BIN_NAME) (debug) to $(CONNECTION)"
	@ssh $(CONNECTION) "systemctl stop $(DEBUG_SYSTEMCTL_UNIT)"
	@ssh $(CONNECTION) "mkdir -p $(DEBUG_DIR)"
	@scp $(DEBUG_BIN_NAME) $(CONNECTION):$(DEBUG_DIR)
	@ssh $(CONNECTION) "chmod +x $(DEBUG_DIR)/$(DEBUG_BIN_NAME)"
	@ssh $(CONNECTION) "systemctl start $(DEBUG_SYSTEMCTL_UNIT)"
	@echo "   Done."

.PHONY: upload-prod
upload-prod:
	@echo ">> Uploading $(PROD_BIN_NAME) (production) to $(CONNECTION)"
	@ssh $(CONNECTION) "systemctl stop $(PROD_SYSTEMCTL_UNIT)"
	@ssh $(CONNECTION) "mkdir -p $(PROD_DIR)"
	@scp $(PROD_BIN_NAME) $(CONNECTION):$(PROD_DIR)
	@ssh $(CONNECTION) "chmod +x $(PROD_DIR)/$(PROD_BIN_NAME)"
	@ssh $(CONNECTION) "systemctl start $(PROD_SYSTEMCTL_UNIT)"
	@echo "   Done."
# END REGION: UPLOADING

# REGION: DEBUG
.PHONY: run-debug
run-debug:
	@echo ">> Running Delve"
	@ssh $(CONNECTION) "systemctl stop $(DEBUG_SYSTEMCTL_UNIT)"
	@ssh $(CONNECTION) "\
		dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient \
		exec $(DEBUG_DIR)/$(DEBUG_BIN_NAME)"
	@ssh $(CONNECTION) "systemctl start $(DEBUG_SYSTEMCTL_UNIT)"

.PHONY: stop-debug
stop-debug:
	@echo ">> Stopping Delve"
	@ssh $(CONNECTION) "pkill dlv || true"

.PHONY: debug
debug: stop-debug build-linux-debug upload-debug run-debug
# END REGION: DEBUG

# REGION: PROD
.PHONY: prod
prod: build-linux-prod upload-prod
# END REGION: PROD

# REGION: WEB
.PHONY: upload-debug-web
upload-debug-web:
	@echo ">> Uploading Web (debug)"
	@cd web/auth && env VITE_API_BASE_URL=$(DEBUG_API_BASE_URL) npm run build
	@ssh $(CONNECTION) "mkdir -p $(DEBUG_WEB_DIR)"
	@scp -r web/auth/dist/* $(CONNECTION):$(DEBUG_WEB_DIR)
	@echo "   Done."

.PHONY: upload-prod-web
upload-prod-web:
	@echo ">> Uploading Web (production)"
	@cd web/auth && env VITE_API_BASE_URL=$(PROD_API_BASE_URL) npm run build
	@ssh $(CONNECTION) "mkdir -p $(PROD_WEB_DIR)"
	@scp -r web/auth/dist/* $(CONNECTION):$(PROD_WEB_DIR)
	@echo "   Done."
# END REGION: WEB
