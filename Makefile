BRANCH = $(CI_BRANCH)
COMMIT = $(CI_COMMIT)
TIMESTAMP = $(shell /bin/date -u +%FT%TZ)
TIMEINT = $(shell /bin/date -u +%Y%m%d%H%M%S)
LD_FLAGS = "-s -X main.branch=${BRANCH} -X main.commit=$(COMMIT) -X main.date=$(TIMESTAMP)"

all: run

clean:
	rm -rf build

build: clean *.go
	CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) -a -installsuffix cgo -o build/csv2api

run: quick
		./build/csv2api

test:
		go test -v ./...

coverage:
	go test . -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

build-linux: clean *.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LD_FLAGS) -a -installsuffix cgo -o build/csv2api.linux

docker: build-linux
	docker build -t mnbbrown/csv2api:latest .