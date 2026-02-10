FROM golang:1.26rc3-trixie AS build-base

WORKDIR /app
COPY . .

RUN go tool github.com/magefile/mage build:release
FROM gcr.io/distroless/static-debian13

WORKDIR /
COPY --from=build-base /app/bin/blog-release /blog

CMD ["./blog"]
