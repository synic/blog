FROM golang:1.23.2-alpine3.20 AS build-base

WORKDIR /app
COPY . .

RUN set -x \
    && go run cmd/compile/*.go -i articles -o ./articles/json -v \
    && go build \
      -tags release \
      -ldflags "-s -w -X main.BuildTime=$(date +%s)" \
      -o bin/blog \
      .

FROM gcr.io/distroless/static-debian12:9efbcaacd8eac4960b315c502adffdbf3398ce62

WORKDIR /
COPY --from=build-base /app/bin/blog /blog

CMD ["./blog", "serve"]
