GOPATH?=`realpath workspace`
BIN="./bin/blog-serve"
CSSDIR=./assets/css
TEMPL_VERSION=v0.2.648
AIR_VERSION=v1.49.0

AIR_TEST := $(shell command -v air 2> /dev/null)
TEMPL_TEST := $(shell command -v templ 2> /dev/null)

.PHONY: dev
dev:
ifndef AIR_TEST
	go install github.com/cosmtrek/air@${AIR_VERSION}
endif
	@air \
	  -root "." \
		-tmp_dir ".tmp" \
		-build.cmd "make build" \
		-build.bin "DEBUG=true ./bin/blog-serve-debug" \
		-build.delay "1000" \
		-build.exclude_dir 'logs,node_modules,bin,assets/css/compiled' \
		-build.exclude_file 'Dockerfile,docker-compose.yaml' \
		-build.exclude_regex '_test.go,.null-ls,_templ.go' \
		-build.include_ext 'go,templ,css,md,js' \
		-build.log "logs/build-errors.log" \
		-misc.clean_on_exit "false"

.PHONY: build
build: clean css templ vet
	@go build -tags debug -o ${BIN}-debug ./cmd/serve

.PHONY: templ
templ:
ifndef TEMPL_TEST
	go install github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}
endif
	@templ generate

.PHONY: release
release: clean-release css templ  vet
	@go build -tags release -ldflags "-s -w" -o ${BIN}-release \
	  ./cmd/serve

.PHONY: clean-common
clean-common:
	@find . -name "*_templ.go" -delete
	@rm -rf ${CSSDIR}/compiled

.PHONY: clean
clean: clean-common
	@rm ${BIN}-debug 2> /dev/null || true

.PHONY: clean-release
clean-release: clean-common
	@rm ${BIN}-release 2> /dev/null || true

.PHONY: test
test:
	@go test -race ./... | tc

.PHONY: vet
vet:
	@go vet ./...

.PHONY: gen-syntax-css
gen-syntax-css:
	@pygmentize -S catppuccin-mocha -f html -a .chroma \
		> ${CSSDIR}/syntax.css

.PHONY: mkdirs
mkdirs:
	@mkdir -p bin
	@mkdir -p ${CSSDIR}/compiled

.PHONY: install-builddeps
install-builddeps:
	npm i

.PHONY: build-container
build-container:
	docker build . -t "blog-serve"

.PHONY: run-container
run-container:
	@docker rm -f blog-serve 2> /dev/null || true
	docker run -p 3000:3000 --name blog-serve -t blog-serve

.PHONY: css
css:
	@./node_modules/.bin/tailwindcss \
		--postcss \
		-i ${CSSDIR}/main.css \
		-o ${CSSDIR}/compiled/main.css \
		--minify
