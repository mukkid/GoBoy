PACKAGE 	= GoBoy
GOPATH		= $(CURDIR)
BASE		= $(GOPATH)/src
TEST		= $(GOPATH)/test
OBJ			= obj
BIN			= bin
SRCFILES	:= $(wildcard $(BASE)/*.go)
OBJECTS 	:= $(patsubst $(BASE)/%.go,$(OBJ)/%.o, $(SRCFILES))

all: | $(BIN)
	go install ./src/goboy
	go build -o $(BIN)/gobjdump ./src/gobjdump
	go get "github.com/stretchr/testify/assert"

.PHONY: test
test: all
	go test ./test
	go test -v ./src/goboy

$(BIN):
	mkdir -p $@

.PHONY: clean
clean:
	rm -rf pkg/ $(BIN)
