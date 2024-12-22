FROM golang:1.24-rc-bookworm AS build-base

WORKDIR /app
COPY . .

RUN go run github.com/magefile/mage@2385abb build:release

FROM gcr.io/distroless/static-debian12:9efbcaacd8eac4960b315c502adffdbf3398ce62

WORKDIR /
COPY --from=build-base /app/bin/blog-release /blog

CMD ["./blog"]
