SHELL := /bin/bash
PKG = github.com/dfuentes/downstream/cmd/downstream
SUBPKGS := \
github.com/dfuentes/downstream/downstream
PKGS = $(PKG) $(SUBPKGS)
VERSION := $(shell cat VERSION)
EXECUTABLE := downstream
BUILDS := \
	build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-linux-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-windows-amd64
COMPRESSED_BUILDS := $(BUILDS:%=%.tar.gz)
RELEASE_ARTIFACTS := $(COMPRESSED_BUILDS:build/%=release/%)

.PHONY: test $(PKGS) clean run

test: $(PKGS)

$(GOPATH)/bin/golint:
	@go get github.com/golang/lint/golint

$(PKGS): $(GOPATH)/bin/golint
	@go get -d -t $@
	@gofmt -w=true $(GOPATH)/src/$@*/**.go
	@echo "LINTING..."
	@$(GOPATH)/bin/golint $(GOPATH)/src/$@*/**.go
	@echo ""
ifeq ($(COVERAGE),1)
	@go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
	@go tool cover -html=$(GOPATH)/src/$@/c.out
else
	@echo "TESTING..."
	@go test $@ -test.v
endif

build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o "$@/$(EXECUTABLE)" $(PKG)
build/$(EXECUTABLE)-v$(VERSION)-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o "$@/$(EXECUTABLE)" $(PKG)
build/$(EXECUTABLE)-v$(VERSION)-windows-amd64:
	GOARCH=amd64 GOOS=windows go build -o "$@/$(EXECUTABLE).exe" $(PKG)
build: $(BUILDS)
%.tar.gz: %
	tar -C `dirname $<` -zcvf "$<.tar.gz" `basename $<`
$(RELEASE_ARTIFACTS): release/% : build/%
	mkdir -p release
	cp $< $@
release: $(RELEASE_ARTIFACTS)

clean:
	rm -rf build release

run:
	@go run cmd/downstream/main.go
