VERSION = 0.0.1

.PHONY: all
build: 
	 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ./bin/$(shell basename `pwd`)-$(VERSION) ./cmd/main.go
