bot:
	go run -v ./cmd/main.go

test:
	go test ./... -cover 

build:
	go build -o pill-reminder ./cmd

# default
run: bot