ifdef update
	u=-u
endif

export GO111MODULE=on

bin/archive:
	go build -ldflags="-s -w" -o ./bin/archive ./cmd/archive

.PHONY: explosion
explosion:
	@cd cmd/explosion && go run main.go tweet.go

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: test
test:
	go test -race ./...