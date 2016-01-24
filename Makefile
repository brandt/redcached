SRCS = $(shell find . -type f -name '*.go')

.PHONY: deps clean

redcached: $(SRCS)
	go build

deps:
	go get -t ./...

clean:
	$(RM) redcached
