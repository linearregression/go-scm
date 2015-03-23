.PHONY: \
	all \
	precommit \
	check_for_codeship \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	generate \
	build \
	install \
	test \
	codeshipsteps \
	cov \
	testcontainer \
	buildcontainer \
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
	go install github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype

updatedeps:
	go get -d -v -u -f ./...
	go install github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype

testdeps: deps
	go get -d -v -t ./...

updatetestdeps: updatedeps
	go get -d -v -t -u -f ./...

generate: deps
	go generate ./...

build: deps generate
	go build ./...

install: deps generate
	go install ./...


test: testdeps generate
	go test -test.v ./...

codeshipsteps: check_for_codeship generate
	codeship steps

cov: testdeps generate
	go get -v github.com/axw/gocov/gocov
	gocov test | gocov report

buildcontainer: testdeps generate
	docker build --file=Dockerfile.build -t pedge/goscmbuild .

testcontainer: testdeps generate
	docker build --file=Dockerfile.test -t pedge/goscmtest .

linuxcompile: deps generate
	bash makebin/compile.sh linux

darwincompile: deps generate
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
