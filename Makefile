.PHONY: \
	all \
	precommit \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	build \
	install \
	test \
	cov \
	libcontainer \
	linuxcompile \
	darwincompile \
	xccompile \
	container \
	linuxdockercompile \
	darwindockercompile \
	xc \
	clean

all: test install

precommit: xc

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

install: deps
	go install ./...

cov: testdeps
	go get -v github.com/axw/gocov/gocov
	gocov test | gocov report

test: testdeps
	go test -test.v ./...

libcontainer: testdeps
	docker build --file=Dockerfile.lib -t pedge/goscmlib .

linuxcompile: deps
	bash makebin/compile.sh linux

darwincompile: deps
	bash makebin/compile.sh darwin

xccompile: linuxcompile darwincompile

container: linuxcompile
	docker build --file=Dockerfile.goscm -t pedge/goscm .

linuxdockercompile: libcontainer
	bash makebin/docker_compile.sh linux

darwindockercompile: libcontainer
	bash makebin/docker_compile.sh darwin

xc: linuxdockercompile darwindockercompile

clean:
	go clean -i ./...
