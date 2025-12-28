run-server:
	go run cmd/server/main.go

run-server-with-mock-readings-ticker:
	MOCK_READINGS_TICKER=true go run cmd/server/main.go

run-client:
	cd cmd/client && npm run dev

.PHONY: run-server run-client run-server-with-mock-readings-ticker
