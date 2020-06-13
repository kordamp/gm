.PHONY: build

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_TIMESTAMP := $(shell date)

build:
	go get ./...
	go build -ldflags "-X 'main.gmVersion=$(GM_VERSION)' -X 'main.gmBuildCommit=$(GIT_COMMIT)' -X 'main.gmBuildTimestamp=$(BUILD_TIMESTAMP)'" gm.go
