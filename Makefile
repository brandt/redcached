SRCS = $(shell find . -type f -name '*.go')

redcached: $(SRCS)
	go build -o redcached cmd/redcached.go

.PHONY: clean
clean:
	$(RM) redcached
