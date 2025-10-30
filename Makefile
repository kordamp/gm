.PHONY: build

build:
	go get ./...
	go build  -o target/${GOOS}-${GOARCH}/ \
      -ldflags "-X 'main.gmVersion=$(GM_VERSION)' -X 'main.gmBuildCommit=$(GIT_COMMIT)' -X 'main.gmBuildTimestamp=$(BUILD_TIMESTAMP)'" gm.go

