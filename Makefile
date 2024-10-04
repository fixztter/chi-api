build:
	@go build -o bin/chi-api cmd/main.go

run: build
	@./bin/chi-api

test:
	@go test -v ./...

clean:
	@rm -rfv ./bin