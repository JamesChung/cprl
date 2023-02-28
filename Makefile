.PHONY: build
build:
	go build -o cprl .

.PHONY: test
test:
	go test -v
