all: build

test: build
	cd goboy && go $@ -v -cover -coverprofile=count.out

build: deps
	cd goboy && go $@
	cd gobjdump && go $@

fmt:
	pushd src/goboy && go fmt && popd
	pushd src/gobjdump && go fmt && popd
	pushd src/gobjdump/main && go fmt && popd

.PHONY: clean
clean:
	cd goboy && go $@ && rm -f count.out
	cd gobjdump && go $@ && rm -f count.out

deps:
	cd goboy && go get -d ./... && go list -f '{{ join .TestImports "\n" }}' | xargs go get -d
	cd gobjdump && go get -d ./... && go list -f '{{ join .TestImports "\n" }}' | xargs go get -d

coverage: test
	cd goboy && sed -i "s/.*\//.\//" count.out && go tool cover -html=count.out
