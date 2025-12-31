APP_NAME := portarbiter
PKG := portarbiter/internal/version

# Git metadata
GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

VERSION := $(subst v,,$(GIT_TAG))

LDFLAGS := \
	-X '$(PKG).Version=$(VERSION)' \
	-X '$(PKG).GitCommit=$(GIT_COMMIT)' \
	-X '$(PKG).BuildDate=$(BUILD_DATE)'

BUILD_DIR := build
DEB_DIR := $(BUILD_DIR)/deb

.PHONY: all build clean deb version

all: build

build:
	@echo "==> Building $(APP_NAME) $(VERSION)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o $(APP_NAME) -ldflags "$(LDFLAGS)" ./cmd/$(APP_NAME)

version:
	@echo "Version:    $(VERSION)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build date: $(BUILD_DATE)"

clean:
	rm -f $(APP_NAME)
	rm -rf $(BUILD_DIR)

deb: build
	@echo "==> Building Debian package"
	rm -rf $(DEB_DIR)
	mkdir -p $(DEB_DIR)/DEBIAN
	mkdir -p $(DEB_DIR)/usr/bin
	mkdir -p $(DEB_DIR)/usr/share/man/man1

	cp $(APP_NAME) $(DEB_DIR)/usr/bin/$(APP_NAME)
	chmod 755 $(DEB_DIR)/usr/bin/$(APP_NAME)

	gzip -c man/$(APP_NAME).1 > $(DEB_DIR)/usr/share/man/man1/$(APP_NAME).1.gz
	chmod 644 $(DEB_DIR)/usr/share/man/man1/$(APP_NAME).1.gz

	sed "s/__VERSION__/$(VERSION)/g" packaging/control > $(DEB_DIR)/DEBIAN/control
	cp packaging/postinst $(DEB_DIR)/DEBIAN/postinst
	chmod 755 $(DEB_DIR)/DEBIAN/postinst

	dpkg-deb --build $(DEB_DIR)
	mv $(DEB_DIR).deb $(APP_NAME)_$(VERSION)_amd64.deb
	@echo "==> Package: $(APP_NAME)_$(VERSION)_amd64.deb"

