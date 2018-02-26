PACKAGE 	= GoBoy
GOPATH		= $(CURDIR)
BASE		= $(GOPATH)/src
TEST		= $(GOPATH)/test
OBJ			= obj
SRCFILES	:= $(wildcard $(BASE)/*.go)
OBJECTS 	:= $(patsubst $(BASE)/%.go,$(OBJ)/%.o, $(SRCFILES))

all:
	go install ./src/goboy

.PHONY: test
test: all
	go test ./test

.PHONY: clean
clean:
	rm -rf pkg/
