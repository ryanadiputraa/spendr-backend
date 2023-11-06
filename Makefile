env:
	cp config/.env.production config/.env.development

server:
	go run cmd/api/main.go development
