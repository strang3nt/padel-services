TARGET_EXEC := tgminiapp
TEMPLATE_DIR := ./template
WEBAPP_DIR := ./cmd/tgminiapp
PKG_DIR := ./pkg
# Added MAIN_DIR or replaced it with WEBAPP_DIR to fix the GO_DIRS variable
MAIN_DIR := $(WEBAPP_DIR)
GO_DIRS := $(TEMPLATE_DIR) $(MAIN_DIR) $(PKG_DIR)
CLIENT_DIR := ./client
CLIENT_BUILD_DIR := $(CLIENT_DIR)/dist

GO_SRCS := $(shell find $(GO_DIRS) -name '*.go') go.mod
CLIENT_SRCS := $(shell find $(CLIENT_DIR) -maxdepth 3 -name '*.ts' -or -name '*.tsx' -or -name '*.html' -or -name '*.json' -or -name '*.css' | grep -v $(CLIENT_BUILD_DIR))

.PHONY: all clean client

all: $(TARGET_EXEC)

# Dedicated target name for the client build
client: $(CLIENT_BUILD_DIR)/index.html

$(TARGET_EXEC): $(GO_SRCS) client
	go build -o $(TARGET_EXEC) github.com/strang3nt/padel-services/cmd/tgminiapp

$(CLIENT_BUILD_DIR)/index.html: $(CLIENT_SRCS)
	cd $(CLIENT_DIR) && pnpm run build
	rm -rf $(WEBAPP_DIR)/dist
	cp -r $(CLIENT_BUILD_DIR) $(WEBAPP_DIR)/dist

clean:
	rm -rf $(TARGET_EXEC) $(CLIENT_DIR)/dist $(WEBAPP_DIR)/dist