.PHONY: build
build:
	go build -o cprl .

.PHONY: test
test:
	go test -v

.PHONY: release
release:
	@git push
	goreleaser release --clean

.PHONY: tag
tag:
	git tag -a v0.0.0-beta.$$(date +"%Y%m%d") -m v0.0.0-beta.$$(date +"%Y%m%d")

.PHONY: local
local:
	go build -o cprl . && mv cprl ~/.local/bin

.PHONY: docs
docs:
	CPRL_DOCS=./docs go run main.go
	@git restore --staged .
	@git add ./docs
	@git commit -m "docs: :memo: update documentation"
