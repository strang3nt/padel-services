TARGET_EXEC := tgminiapp

TEMPLATE_DIR := ./template
WEBAPP_DIR := ./cmd/tgminiapp
PKG_DIR := ./pkg
GO_DIRS := $(TEMPLATE_DIR) $(MAIN_DIR) $(PKG_DIR)
CLIENT_DIR := ./client
CLIENT_BUILD_DIR := $(CLIENT_DIR)/dist

# Go source files
GO_SRCS := $(shell find $(GO_DIRS) -name '*.go') go.mod

# client source files, i.e. typescript, react-ts, and configuration files
CLIENT_SRCS := $(shell find $(CLIENT_DIR) -path $(CLIENT_BUILD_DIR) -prune -name '*.ts' -or -name '*.tsx' -or -name '*.html' -or -name '*.json' -or -name '*.css' -not -path $(CLIENT_DIR)/dist)

.PHONY: all clean

all: $(TARGET_EXEC)

# Vite generated webapp files
CLIENT_OBJS := $(CLIENT_BUILD_DIR)/index.html

$(TARGET_EXEC): $(GO_SRCS) $(CLIENT_OBJS)
	go build github.com/strang3nt/padel-services/cmd/tgminiapp

$(CLIENT_OBJS): $(CLIENT_SRCS)
	cd client && npm run build:dev
	rm -rf $(WEBAPP_DIR)/dist
	cp -r $(CLIENT_BUILD_DIR) $(WEBAPP_DIR)/dist

clean:
	rm -rf $(TARGET_EXEC) $(CLIENT_DIR)/dist $(WEBAPP_DIR)/dist