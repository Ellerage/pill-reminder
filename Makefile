bot:
	go run -v ./cmd/main.go

test:
	go test ./... -cover 

build_prod:
	GOOS=linux GOARCH=amd64 go build -o pill-reminder ./cmd

# default
run: bot