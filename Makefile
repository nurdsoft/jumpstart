NAME = jumpstart

BUILDTIME = $(shell date '+%s')

TAG=$(shell git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
COMMIT=$(shell git rev-parse --short HEAD)

# VERSION is the latest tag, or branch name if no tags
VERSION=$(if $(TAG),$(TAG),$(BRANCH))

PKG_PATH = github.com/nurdsoft/$(NAME)/pkg

RELEASES_DIR = releases
BUILD_OPTS = -ldflags "-s -w -X $(PKG_PATH).VERSION=$(VERSION) -X $(PKG_PATH).COMMIT=$(COMMIT) -X $(PKG_PATH).BUILDTIME=$(BUILDTIME)"

TEMP_FILE = tmp

clean:
	rm -rf ./$(TEMP_FILE)
	rm -rf ./$(RELEASES_DIR)

$(RELEASES_DIR):
	mkdir -p $(RELEASES_DIR)

$(RELEASES_DIR)/$(NAME): $(RELEASES_DIR)
	CGO_ENABLED=0 go build -o $@ $(BUILD_OPTS) ./cmd/	

$(RELEASES_DIR)/$(NAME)-%: $(RELEASES_DIR)
	GOOS=$* CGO_ENABLED=0 GOARCH=amd64 go build -o $@ $(BUILD_OPTS) ./cmd/

# .PHONY: docker
# docker
# 	docker build . -t $(NAME) -t $(NAME):$(VERSION)-$(COMMIT) -t $(NAME):$(VERSION)

.PHONY: all
all: linux_amd64 linux_arm64 linux_arm darwin_amd64 darwin_arm64 windows_amd64 windows_arm64 windows_arm

linux_arm64:
	GOOS=linux GOARCH=arm64 go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-linux-arm64 ./cmd/

linux_amd64:
	GOOS=linux GOARCH=amd64 go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-linux-amd64 ./cmd/

linux_arm:
	GOOS=linux GOARCH=arm go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-linux-arm ./cmd/

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-darwin-arm64 ./cmd/

darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-darwin-amd64 ./cmd/

windows_amd64:
	GOOS=windows GOARCH=amd64 go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-windows-amd64.exe ./cmd/

windows_arm64:
	GOOS=windows GOARCH=arm64 go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-windows-arm64.exe ./cmd/

windows_arm:
	GOOS=windows GOARCH=arm go build $(BUILD_OPTS) -o $(RELEASES_DIR)/$(NAME)-windows-arm.exe ./cmd/