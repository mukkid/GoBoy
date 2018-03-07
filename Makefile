PACKAGE 	= GoBoy
GOPATH		= $(CURDIR)
BASE		= $(GOPATH)/src
TEST		= $(GOPATH)/test
OBJ			= obj
BIN			= bin
SRCFILES	:= $(wildcard $(BASE)/*.go)
OBJECTS 	:= $(patsubst $(BASE)/%.go,$(OBJ)/%.o, $(SRCFILES))

all: | $(BIN)
	go get "github.com/stretchr/testify/assert"
	go get -u github.com/go-gl/glfw/v3.2/glfw
	go get "github.com/hajimehoshi/ebiten"
	go build -o $(BIN)/goboy ./src/goboy
	go build -o $(BIN)/gobjdump ./src/gobjdump

.PHONY: test
test: all
	go test ./test
	go test -v ./src/goboy

$(BIN):
	mkdir -p $@

.PHONY: clean
clean:
	rm -rf pkg/ $(BIN)
