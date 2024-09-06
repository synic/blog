GOPATH?=`realpath workspace`
BIN=./bin/blog
TEMPL_VERSION=v0.2.663
AIR_VERSION=v1.49.0

AIR_TEST := $(shell command -v air 2> /dev/null)
TEMPL_TEST := $(shell command -v templ 2> /dev/null)
TAILWIND_TEST := $(shell command -v ./node_modules/.bin/tailwindcss 2> /dev/null)

.PHONY: dev
dev: install-devdeps
	@air \
	  -root "." \
		-tmp_dir ".tmp" \
		-build.cmd "make build" \
		-build.bin "DEBUG=true ./bin/blog-debug serve" \
		-build.delay "1000" \
		-build.exclude_dir 'logs,node_modules,bin,assets/css/compiled' \
		-build.exclude_file 'Dockerfile,docker-compose.yaml' \
		-build.exclude_regex '_test.go,.null-ls,_templ.go' \
		-build.include_ext 'go,templ,css,md,js' \
		-build.log "logs/build-errors.log" \
		-misc.clean_on_exit "false"

.PHONY: install-devdeps
install-devdeps:
ifndef AIR_TEST
	go install github.com/cosmtrek/air@${AIR_VERSION}
endif
ifndef TEMPL_TEST
	go install github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}
endif
ifndef TAILWIND_TEST
	npm i
endif

.PHONY: build
build: codegen
	go build -tags debug -o ${BIN}-debug .

.PHONY: codegen
codegen:
	go generate ./...

.PHONY: release
release: vet
	@go build -tags release -ldflags "-s -w" -o ${BIN}-release .

.PHONY: test
test:
	go test -race ./... | tc

.PHONY: vet
vet:
	go vet ./...

.PHONY: gen-syntax-css
gen-syntax-css:
	pygmentize -S catppuccin-mocha -f html -a .chroma \
		> ${CSSDIR}/syntax.css

.PHONY: mkdirs
mkdirs:
	mkdir -p bin

.PHONY: build-container
build-container:
	docker build . -t "blog-serve"

.PHONY: run-container
run-container:
	docker rm -f blog-serve 2> /dev/null || true
	docker run -p 3000:3000 --name blog-serve -t blog-serve
