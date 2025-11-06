.PHONY: build

.EXPORT_ALL_VARIABLES:
CGO_ENABLED=0

build:
	go get ./...
	go build -o target/${GOOS}-${GOARCH}/ \
      -trimpath \
      -ldflags "-s -w -buildid= -X 'main.gmVersion=$(GM_VERSION)' -X 'main.gmBuildCommit=$(GIT_COMMIT)' -X 'main.gmBuildTimestamp=$(BUILD_TIMESTAMP)'" gm.go

