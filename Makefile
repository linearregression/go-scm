.PHONY: all deps updatedepst testdeps updatetestdeps build test cov install compile container doc clean

all: test install

deps:
	go get -d -v ./...

updatedeps:
	go get -d -v -u -f ./...

testdeps: deps
	go get -d -v -t ./...

updatetestdeps: updatedeps
	go get -d -v -t -u -f ./...

build: deps
	go build ./...

test: testdeps
	go test -test.v ./...

cov: testdeps
	go get -v github.com/axw/gocov/gocov
	gocov test | gocov report

install: deps
	go install ./...

compile: deps
	mkdir -p tmp
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-s' -o tmp/go-scm cmd/go-scm/main.go
	ls -lh tmp

container: compile
	docker build -t pedge/goscm .

doc:
	go get -v github.com/robertkrimen/godocdown/godocdown
	cp .readme.header README.md
	godocdown | tail -n +7 >> README.md

clean:
	go clean -i ./...
	docker rmi pedge/goscm || true
