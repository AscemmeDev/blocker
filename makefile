run:
	 go run ./cmd/main.go

build:
	 go build -mod=vendor -ldflags="-w -s" -o blocker ./cmd