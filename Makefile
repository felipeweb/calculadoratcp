COMMIT_HASH=`git rev-parse --short HEAD>/dev/null`
DIST := dist
IMPORT := github.com/felipeweb/calculadoratcp
OUT := 'calculator'
TARGETS ?= linux/*,darwin/*,windows/*

EXECUTABLE := $(OUT)
ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(OUT).exe
endif

PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)

.PHONY: errcheck
errcheck:
	@which errcheck > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

.PHONY: install
install: $(wildcard *.go)
	go install

.PHONY: build
build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	go build -v -ldflags '-s -w $(LDFLAGS)' -o $@

.PHONY: release
release: release-dirs release-build release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries

.PHONY: release-build
release-build:
	@which xgo > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/karalabe/xgo; \
	fi
	xgo -dest $(DIST)/binaries -targets '$(TARGETS)' -out $(EXECUTABLE) $(IMPORT)
ifeq ($(CI),drone)
	mv /build/* $(DIST)/binaries
endif

.PHONY: release-check
release-check:
	cd $(DIST)/binaries; $(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)
