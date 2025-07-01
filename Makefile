bot:
	go run -v ./cmd/main.go

test:
	go test ./... -cover 

# default
run: bot