PACKAGE 	= GoBoy
GOPATH		= $(CURDIR)
BASE		= $(GOPATH)/src
TEST		= $(GOPATH)/test

.PHONY: test
test:
	go test ./test
