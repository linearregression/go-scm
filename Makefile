.PHONY: \
	all \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	generate \
	build \
	install \
	lint \
	vet \
	errcheck \
	pretest \
	test \
	testlong \
	clean

all: test install

deps:
	go get -d -v ./...
	go get -v github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype

updatedeps:
	go get -d -v -u -f ./...
	go get -v -u -f github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype

testdeps:
	go get -d -v -t ./...
	go get -v github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype
	go get -v golang.org/x/tools/cmd/vet
	go get -v github.com/kisielk/errcheck
	go get -v github.com/golang/lint/golint

updatetestdeps:
	go get -d -v -t -u -f ./...
	go get -v -u -f github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype
	go get -v -u -f golang.org/x/tools/cmd/vet
	go get -v -u -f github.com/kisielk/errcheck
	go get -v -u -f github.com/golang/lint/golint

generate:
	go generate ./...

build: deps generate
	go build ./...

install: deps generate
	go install ./...

lint: testdeps generate
	-golint ./...

vet: testdeps generate
	go vet ./...

errcheck: testdeps generate
	errcheck ./...

pretest: lint vet errcheck

test: pretest
	go test -test.v -test.short ./...

testlong: pretest
	go test -test.v ./...

clean:
	go clean -i ./...
