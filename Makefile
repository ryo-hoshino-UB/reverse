GO_CMD=go

fmt:
	cd api/app && $(GO_CMD) fmt ./...

sqlc:
	cd api/app && sqlc generate