.PHONY: \
	all \
	precommit \
	check_for_codeship \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	build \
	install \
	test \
	codeshipsteps \
	cov \
	libcontainer \
	linuxcompile \
	darwincompile \
	xccompile \
	container \
	linuxdockercompile \
	darwindockercompile \
	doc \
	xc \
	clean

all: test install

precommit: xc

check_for_codeship:
	@ if ! which codeship > /dev/null; then \
			echo "error: codeship not installed" >&2; \
	  fi
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


test: testdeps
	go test -test.v ./...

codeshipsteps: check_for_codeship 
	codeship steps

cov: testdeps
	go get -v github.com/axw/gocov/gocov
	gocov test | gocov report

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

doc:
	go get -v github.com/robertkrimen/godocdown/godocdown
	cp .readme.header README.md
	godocdown | tail -n +7 >> README.md

xc: linuxdockercompile darwindockercompile

clean:
	go clean -i ./...
