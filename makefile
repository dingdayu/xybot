.PHONY: clean linux fmt

##export GOPATH:=$(shell pwd)

linux :
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


clean :
	rm -rf wxbot