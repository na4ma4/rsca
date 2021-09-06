GO_MATRIX += darwin/amd64
GO_MATRIX += linux/amd64
GO_MATRIX += windows/amd64

APP_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH ?= $(shell git show -s --format=%h)

GO_DEBUG_ARGS   ?= -v -ldflags "-X main.version=$(GO_APP_VERSION)+debug -X main.commit=$(GIT_HASH) -X main.date=$(APP_DATE) -X main.builtBy=makefiles"
GO_RELEASE_ARGS ?= -v -ldflags "-X main.version=$(GO_APP_VERSION) -X main.commit=$(GIT_HASH) -X main.date=$(APP_DATE) -X main.builtBy=makefiles -s -w"

GENERATED_FILES += artifacts/certs/ca.pem
GENERATED_FILES += artifacts/certs/server.pem
GENERATED_FILES += artifacts/certs/server-key.pem
GENERATED_FILES += artifacts/certs/client.pem
GENERATED_FILES += artifacts/certs/client-key.pem
GENERATED_FILES += artifacts/certs/cert.pem
GENERATED_FILES += artifacts/certs/key.pem

GO_TEST_REQ += test/test.cmd

-include .makefiles/Makefile
-include .makefiles/pkg/protobuf/v2/Makefile
-include .makefiles/pkg/go/v1/Makefile
-include .makefiles/ext/na4ma4/lib/golangci-lint/v1/Makefile
-include .makefiles/ext/na4ma4/lib/cfssl/v1/Makefile
-include .makefiles/ext/na4ma4/lib/goreleaser/v1/Makefile

.makefiles/ext/na4ma4/%: .makefiles/Makefile
	@curl -sfL https://raw.githubusercontent.com/na4ma4/makefiles-ext/main/v1/install | bash /dev/stdin "$@"

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


######################
# Linting
######################

ci:: lint
