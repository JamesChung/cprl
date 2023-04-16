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
tag: update docs
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

.PHONY: update
update:
	go get -u
	go mod tidy
	@git restore --staged .
	@git add go.mod go.sum
	@git commit -m "chore: :arrow_up: upgrade dependencies"
