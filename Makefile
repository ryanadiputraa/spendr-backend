env:
	cp config/config.example.yml config/config.yml

server:
	go run cmd/api/main.go development
