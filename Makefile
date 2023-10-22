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