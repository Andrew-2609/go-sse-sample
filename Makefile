run-server:
	go run cmd/server/main.go

run-client:
	node cmd/client/client.mjs

.PHONY: run-server run-client
