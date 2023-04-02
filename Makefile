.PHONY: build
build:
	go build -o cprl .

.PHONY: test
test:
	go test -v

release:
	goreleaser release --clean

.PHONY: local
local:
	go build -o cprl . && mv cprl ~/.local/bin
