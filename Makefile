GO_MATRIX_OS ?= darwin linux windows
GO_MATRIX_ARCH ?= amd64

APP_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH ?= $(shell git show -s --format=%h)

GO_DEBUG_ARGS   ?= -v -ldflags "-X main.version=$(GO_APP_VERSION)+debug -X main.gitHash=$(GIT_HASH) -X main.buildDate=$(APP_DATE)"
GO_RELEASE_ARGS ?= -v -ldflags "-X main.version=$(GO_APP_VERSION) -X main.gitHash=$(GIT_HASH) -X main.buildDate=$(APP_DATE) -s -w"

_GO_GTE_1_14 := $(shell expr `go version | cut -d' ' -f 3 | tr -d 'a-z' | cut -d'.' -f2` \>= 14)
ifeq "$(_GO_GTE_1_14)" "1"
_MODFILEARG := -modfile tools.mod
endif

GENERATED_FILES += artifacts/certs/ca.pem
GENERATED_FILES += artifacts/certs/server.pem
GENERATED_FILES += artifacts/certs/server-key.pem
GENERATED_FILES += artifacts/certs/client.pem
GENERATED_FILES += artifacts/certs/client-key.pem
GENERATED_FILES += artifacts/certs/cert.pem
GENERATED_FILES += artifacts/certs/key.pem
GENERATED_FILES += test/test.cmd

-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v1/Makefile
-include .makefiles/pkg/go/v1/Makefile

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

.PHONY: run
run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/rsca
	"$<" $(RUN_ARGS)

.PHONY: run-admin
run-admin: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/rsc
	"$<" $(RUN_ARGS)

.PHONY: run-daemon
run-daemon: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/rscad
	"$<" $(RUN_ARGS)

.PHONY: install
install: $(REQ) $(_SRC) | $(USE)
	$(eval PARTS := $(subst /, ,$*))
	$(eval BUILD := $(word 1,$(PARTS)))
	$(eval OS    := $(word 2,$(PARTS)))
	$(eval ARCH  := $(word 3,$(PARTS)))
	$(eval BIN   := $(word 4,$(PARTS)))
	$(eval ARGS  := $(if $(findstring debug,$(BUILD)),$(DEBUG_ARGS),$(RELEASE_ARGS)))

	CGO_ENABLED=$(CGO_ENABLED) GOOS="$(OS)" GOARCH="$(ARCH)" go install $(ARGS) "./cmd/..."


######################
# Custom
######################

artifacts/protobuf/go.proto_paths.jq: artifacts/protobuf/go.proto_paths
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	jq -Rn 'inputs | select(.)' < "$(^)" > "$(@)"

.vscode/settings.json: artifacts/protobuf/go.proto_paths.jq
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	$(if $(shell cat "$(@)" 2>/dev/null),cat "$(@)",echo '{}') | jq --slurpfile po "$(<)" '.protoc.options=$$po' > "$(@).tmp"
	mv "$(@).tmp" "$(@)"


######################
# Test
######################

test/test.cmd:
	ln -s /dev/null "$(@)"


# .PHONY: test-feed-server
# test-feed-server: test/feed/Dockerfile test/feed/gotime.xml
# 	-@docker stop archpod-test-server
# 	docker build -t archpod-test-server:latest "$(<D)"
# 	docker run --rm -d --name "archpod-test-server" -p 8018:80/tcp archpod-test-server:latest


######################
# CFSSL
######################

CFSSL := artifacts/bin/cfssl
$(CFSSL):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) github.com/cloudflare/cfssl/cmd/cfssl

CFSSLJSON := artifacts/bin/cfssljson
$(CFSSLJSON):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) github.com/cloudflare/cfssl/cmd/cfssljson

.PHONY: cfssl
cfssl: artifacts/certs/server.pem artifacts/certs/client.pem

artifacts/certs/ca-config.json: test/ca-config.json
	-@mkdir -p "$(@D)"
	cp "$(<)" "$(@)"

artifacts/certs/ca.pem: $(CFSSL) $(CFSSLJSON) artifacts/certs/ca-config.json test/ca-csr.json
	-@mkdir -p "$(@D)"
	$(CFSSL) gencert -initca -config="artifacts/certs/ca-config.json" -profile="ca" test/ca-csr.json | $(CFSSLJSON) -bare artifacts/certs/ca -
	$(CFSSL) sign -ca="artifacts/certs/ca.pem" -ca-key="artifacts/certs/ca-key.pem" -config="artifacts/certs/ca-config.json" -profile="ca" -csr=artifacts/certs/ca.csr test/ca-csr.json | $(CFSSLJSON) -bare artifacts/certs/ca

artifacts/certs/cert.pem: test/admin.json $(CFSSL) $(CFSSLJSON) artifacts/certs/ca.pem
	-@mkdir -p "$(@D)"
	$(CFSSL) gencert -initca -config="artifacts/certs/ca-config.json" -profile="client" "$(<)" | $(CFSSLJSON) -bare artifacts/certs/cert -
	$(CFSSL) sign -ca="artifacts/certs/ca.pem" -ca-key="artifacts/certs/ca-key.pem" -config="artifacts/certs/ca-config.json" -profile="client" artifacts/certs/cert.csr | $(CFSSLJSON) -bare artifacts/certs/cert

artifacts/certs/key.pem: artifacts/certs/cert.pem
	-@mkdir -p "$(@D)"
	cp artifacts/certs/cert-key.pem "$(@)"

artifacts/certs/server.pem: test/host.json $(CFSSL) $(CFSSLJSON) artifacts/certs/ca.pem
	-@mkdir -p "$(@D)"
	$(CFSSL) gencert -initca -config="artifacts/certs/ca-config.json" -profile="server" "$(<)" | $(CFSSLJSON) -bare artifacts/certs/server -
	$(CFSSL) sign -ca="artifacts/certs/ca.pem" -ca-key="artifacts/certs/ca-key.pem" -config="artifacts/certs/ca-config.json" -profile="server" artifacts/certs/server.csr | $(CFSSLJSON) -bare artifacts/certs/server

artifacts/certs/client.pem: test/client.json $(CFSSL) $(CFSSLJSON) artifacts/certs/ca.pem
	-@mkdir -p "$(@D)"
	$(CFSSL) gencert -initca -config="artifacts/certs/ca-config.json" -profile="client" "$(<)" | $(CFSSLJSON) -bare artifacts/certs/client -
	$(CFSSL) sign -ca="artifacts/certs/ca.pem" -ca-key="artifacts/certs/ca-key.pem" -config="artifacts/certs/ca-config.json" -profile="client" artifacts/certs/client.csr | $(CFSSLJSON) -bare artifacts/certs/client


######################
# Linting
######################

MISSPELL := artifacts/bin/misspell
$(MISSPELL):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) github.com/client9/misspell/cmd/misspell

GOLINT := artifacts/bin/golint
$(GOLINT):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) golang.org/x/lint/golint

GOLANGCILINT := artifacts/bin/golangci-lint
$(GOLANGCILINT):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(MF_PROJECT_ROOT)/$(@D)" v1.33.0

STATICCHECK := artifacts/bin/staticcheck
$(STATICCHECK):
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	GOBIN="$(MF_PROJECT_ROOT)/$(@D)" go get $(_MODFILEARG) honnef.co/go/tools/cmd/staticcheck

artifacts/cover/staticheck/unused-graph.txt: $(STATICCHECK) $(GO_SOURCE_FILES)
	-@mkdir -p "$(MF_PROJECT_ROOT)/$(@D)"
	$(STATICCHECK) -debug.unused-graph "$(@)" ./...
	# cat "$(@)"

.PHONY: lint
lint:: $(GOLINT) $(MISSPELL) $(GOLANGCILINT) $(STATICCHECK) artifacts/cover/staticheck/unused-graph.txt
	go vet ./...
	$(GOLINT) -set_exit_status ./...
	$(MISSPELL) -w -error -locale UK ./...
	$(GOLANGCILINT) run --enable-all --disable 'exhaustivestruct,paralleltest' ./...
	$(STATICCHECK) -fail "all,-U1001" ./...

ci:: lint


######################
# Preload Tools
######################

.PHONY: tools
tools: $(MISSPELL) $(GOLINT) $(GOLANGCILINT) $(STATICCHECK) $(CFSSL) $(CFSSLJSON)
