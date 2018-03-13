all: build

test: build
	cd goboy && go $@ -v
	cd gobjdump && go $@ -v

build: deps
	cd goboy && go $@
	cd gobjdump && go $@

clean:
	cd goboy && go $@
	cd gobjdump && go $@

deps:
	cd goboy && go get -d ./... && go list -f '{{ join .TestImports "\n" }}' | xargs go get -d
	cd gobjdump && go get -d ./... && go list -f '{{ join .TestImports "\n" }}' | xargs go get -d
