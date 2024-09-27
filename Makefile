GOPATH?=`realpath workspace`
BIN=./bin/blog
TEMPL_VERSION=v0.2.778
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
		-build.exclude_dir 'logs,node_modules,bin' \
		-build.exclude_file 'Dockerfile,docker-compose.yaml' \
		-build.exclude_regex '_test.go,.null-ls,_templ.go' \
		-build.include_ext 'go,templ,css,json,js' \
		-build.stop_on_error "true"
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
build: codegen vet
	go build -race -tags debug -o ${BIN}-debug ./cmd/serve/serve.go

.PHONY: codegen
codegen: install-devdeps
	templ generate
	@./node_modules/.bin/tailwindcss --postcss \
		-i ./internal/web/css/main.css \
		-o ./cmd/serve/assets/css/main.css --minify

.PHONY: articles
articles:
	@go run cmd/compile/compile.go -i ./articles -o ./cmd/serve/articles

.PHONY: articles-recompile
articles-recompile:
	@go run cmd/compile/compile.go -i ./articles -o ./cmd/serve/articles -recompile -v

.PHONY: release
release:
	go build -tags release -ldflags "-s -w" -o ${BIN}-release ./cmd/serve/serve.go

.PHONY: test
test:
	go test -race ./...

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
	docker build --no-cache . -t "blog-serve"

.PHONY: run-container
run-container:
	docker rm -f blog-serve 2> /dev/null || true
	docker run -p 3000:3000 --name blog-serve -t blog-serve
