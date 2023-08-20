BINNAME=main
MODULE_NAME=web-api-searching-and-pagination

default: air

run: build
	./bin/${BINNAME}

build:
	 go build -race -ldflags="-X '$(MODULE_NAME)/src/link-constants.debug=true' -X '$(MODULE_NAME)/src/link-constants.version=0.0.1'" -o bin/${BINNAME} src/main.go

release:
	go build -o bin/$(BINNAME) -ldflags="-s -X '$(MODULE_NAME)/src/link-constants.version=0.0.1'" src/main.go

test:
	go test ./...

air:
	air -build.cmd "make build" -build.bin "bin/main"
