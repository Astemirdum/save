
run:
	docker-compose up --build

run-client:
	cd client/cmd; go run main.go;

run-server:
	cd server/cmd; go run main.go;

lint:
	golangci-lint run --fix
	gofumpt -w -s ./..

test:
	 go test -v ./...

mocks:
	cd server/internal/service/mocks/; go generate;

.PHONY: run, run-client, run-server, lint, test, mocks