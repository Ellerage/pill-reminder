bot:
	go run -v ./cmd/main.go

test:
	go test ./... -cover

test-i:
	go test -v ./tests/...

build_prod:
	GOOS=linux GOARCH=amd64 go build -o pill-reminder ./cmd

lint:
	golangci-lint run --fix

# default
run: bot