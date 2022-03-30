export CGO_ENABLED=0
 
export VERBOSE = @

ifeq "$(V)" "1"
	VERBOSE =
endif

default: build
.PHONY: default clean test lint

PROG=server
PROG_SRC=cmd/http-server/main.go

PROG_BUILD_CGO_FLAGS?=-ldflags "-X main.HashCommit=`git rev-parse HEAD` -X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`"
PROG_BUILD_GO_FLAGS?=

PROG_BUILD_FLAGS?=${PROG_BUILD_GO_FLAGS} ${PROG_BUILD_CGO_FLAGS}

${PROG}: ${PROG_SRC}
	$(VERBOSE) echo "-> build binary ..."
	$(VERBOSE) go build ${PROG_BUILD_FLAGS} -o $@ $^

clean:
	rm ${PROG}

build: clean ${PROG}

test:
	$(VERBOSE) echo "-> running tests ..."
	$(VERBOSE) CGO_ENABLED=1 go test -race ./...
 
lint:
	$(VERBOSE) echo "-> running linters ..."
	$(VERBOSE) golangci-lint run ./...
