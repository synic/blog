FROM golang:1.23.2-alpine3.20 AS build-base

WORKDIR /app

COPY . .

ENV GOPATH=/go
ENV PATH="${PATH}:/go/bin"
RUN set -x \
    && go run cmd/compile/compile.go -i articles -o cmd/serve/articles -v \
    && go build -tags release -ldflags "-s -w" -o bin/blog cmd/serve/serve.go

FROM gcr.io/distroless/static-debian12:9efbcaacd8eac4960b315c502adffdbf3398ce62

WORKDIR /
COPY --from=build-base /app/bin/blog /blog

CMD ["./blog", "serve"]
