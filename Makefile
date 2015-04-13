.PHONY: \
	all \
	precommit \
	check_for_jet \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	generate \
	build \
	install \
	test \
	jetsteps \
	cov \
	testcontainer \
	buildcontainer \
	clonecontainer \
	tarballcontainer \
	linuxcompile \
	darwincompile \
	xccompile \
	linuxdockercompile \
	darwindockercompile \
	doc \
	xc \
	clean

all: test install

precommit: xc

check_for_jet:
	@ if ! which jet > /dev/null; then \
			echo "error: jet not installed" >&2; \
			exit 1; \
	  fi

deps:
	go get -d -v ./...
	go get github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype
	go get github.com/peter-edge/go-record/cmd/gen-record

updatedeps:
	go get -d -v -u -f ./...
	go get github.com/peter-edge/go-gen-enumtype/cmd/gen-enumtype
	go get github.com/peter-edge/go-record/cmd/gen-record

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

jetsteps: check_for_jet generate
	jet steps

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

linuxdockercompile: buildcontainer
	bash makebin/docker_compile.sh linux

darwindockercompile: buildcontainer
	bash makebin/docker_compile.sh darwin

clonecontainer: linuxdockercompile
	docker build --file=Dockerfile.clone -t pedge/goscmclone .

tarballcontainer: linuxdockercompile
	docker build --file=Dockerfile.tarball -t pedge/goscmtarball .

doc:
	go get -v github.com/robertkrimen/godocdown/godocdown
	cp .readme.header README.md
	godocdown | tail -n +7 >> README.md

xc: linuxdockercompile darwindockercompile

clean:
	go clean -i ./...
