FROM golang:1.23.2-alpine3.20 AS build-base

WORKDIR /app
COPY . .

RUN set -x \
    && go run cmd/compile/*.go -i articles -o ./articles/json -v \
    && go run github.com/magefile/mage@2385abb build:release

FROM gcr.io/distroless/static-debian12:9efbcaacd8eac4960b315c502adffdbf3398ce62

WORKDIR /
COPY --from=build-base /app/bin/blog-release /blog

CMD ["./blog", "serve"]
