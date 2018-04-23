
build_linux=GOOS=linux GOARCH=amd64 go build
build=go build


.PHONY: build build_linux

clean:
	go clean ./...
	-rm -vr build
	-rm -v restapi dbstore


restapi:
	$(build) cmd/restapi/restapi.go

dbstore:
	$(build) cmd/dbstore/dbstore.go

build: restapi dbstore

restapi_linux:
	$(build_linux) -o build/linux/restapi cmd/restapi/restapi.go

dbstore_linux:
	$(build_linux) -o build/linux/dbstore cmd/dbstore/dbstore.go

build_linux: restapi_linux dbstore_linux

docker: clean build_linux
	docker build -t urllookup -f Dockerfile .
