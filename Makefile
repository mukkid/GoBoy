SRCS = goboy gobjdump

all: build

test: build
	cd goboy && go $@ -v -cover -coverprofile=count.out

build: $(SRCS)

$(SRCS): deps 
	cd $@ && go build

fmt:
	pushd goboy && go fmt && popd
	pushd gobjdump && go fmt && popd

.PHONY: clean
clean:
	cd goboy && go $@ && rm -f count.out
	cd gobjdump && go $@ && rm -f count.out

deps:
	cd goboy && go get -d ./... && go list -f '{{ join .TestImports "\n" }}' | xargs go get -d
	cd gobjdump && go get -d ./... && go list -f '{{ join .TestImports "\n" }}' | xargs go get -d

coverage: test
	cd goboy && sed -i "s/.*\//.\//" count.out && go tool cover -html=count.out
