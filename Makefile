.PHONY: build
build:
	go build -o cprl .

.PHONY: test
test:
	go test -v

.PHONY:release
release:
	goreleaser release --clean

.PHONY: tag
tag:
	git tag -a v0.0.0-alpha-$$(date +"%Y%m%d%H%M%S") -m v0.0.0-alpha-$$(date +"%Y%m%d%H%M%S")

.PHONY: local
local:
	go build -o cprl . && mv cprl ~/.local/bin

.PHONY: docs
docs:
	CPRL_DOCS=./docs go run main.go
