PACKAGE 	= GoBoy
GOPATH		= $(CURDIR)
BASE		= $(GOPATH)/src
TEST		= $(GOPATH)/test
OBJ			= obj
BIN			= bin
SRCFILES	:= $(wildcard $(BASE)/*.go)
OBJECTS 	:= $(patsubst $(BASE)/%.go,$(OBJ)/%.o, $(SRCFILES))

all: fmt | $(BIN)
	go get "github.com/stretchr/testify/assert"
	go get -u github.com/go-gl/glfw/v3.2/glfw
	go get "github.com/hajimehoshi/ebiten"
	go build -o $(BIN)/goboy ./src/goboy
	go get "github.com/pborman/getopt/v2"
	go install ./src/gobjdump
	go build -o $(BIN)/gobjdump ./src/gobjdump/main

.PHONY: test
test: all
	go test ./test
	go test -v ./src/goboy

$(BIN):
	mkdir -p $@

fmt:
	pushd src/goboy && go fmt && popd
	pushd src/gobjdump && go fmt && popd
	pushd src/gobjdump/main && go fmt && popd

.PHONY: clean
clean:
	rm -rf pkg/ $(BIN)
