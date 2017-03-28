
APP=consul2file

all:
	docker run --rm -v "${PWD}":/go/src/${APP} -w /go/src/${APP} golang sh -c "go get -v -d .; go build -o /go/src/${APP}/bin/${APP}"