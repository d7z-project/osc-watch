GOPATH := $(shell go env GOPATH)

GOARCH ?= $(shell go env GOARCH)

GOOS ?= $(shell go env GOOS)

.PHONY: fmt all clean test binary


all:
	@make binary GOOS=linux GOARCH=amd64 && make binary GOOS=linux GOARCH=arm64 && \
	 make binary GOOS=windows GOARCH=amd64

clean:
	@rm -rf osc-watch build osc-watch*


EXEC_EXT = 
ifeq ($(GOOS),windows)
EXEC_EXT = .exe
endif

binary: osc-watch-$(GOOS)-$(GOARCH)$(EXEC_EXT).gz


osc-watch-$(GOOS)-$(GOARCH)$(EXEC_EXT): $(shell find ./ -type f -name '*.go') go.mod go.sum fmt clean
	 @GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build  -ldflags "-w -s" -v -o $@

osc-watch-$(GOOS)-$(GOARCH)$(EXEC_EXT).gz: osc-watch-$(GOOS)-$(GOARCH)$(EXEC_EXT)
	@gzip osc-watch-$(GOOS)-$(GOARCH)$(EXEC_EXT) -9

fmt:
	@(test -f "$(GOPATH)/bin/gofumpt" || go install mvdan.cc/gofumpt@latest) && \
	"$(GOPATH)/bin/gofumpt" -l -w .

test: osc-watch
	@go test -v ./...