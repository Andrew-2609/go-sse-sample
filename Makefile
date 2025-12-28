run-server:
	go run cmd/server/main.go

run-client:
	cd cmd/client && npm run dev

.PHONY: run-server run-client
