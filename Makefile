CMD = taskmasterd taskmasterctl tmtest
CMDDIR = ./cmd

.PHONY: build
build: $(CMD)

$(CMD): %: $(CMDDIR)/%/main.go
	go build -o $@ $<

.PHONY: clean
clean:
	rm -f $(CMD)

.PHONY: re
re: clean build

.PHONY: test
test:
	go test ./...

.PHONY: genproto
genproto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		./api/*/*.proto