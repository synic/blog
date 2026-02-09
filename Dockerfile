FROM golang:1.26rc3-trixie AS build-base

WORKDIR /app
COPY . .

RUN go run github.com/magefile/mage@2385abb build:release
FROM gcr.io/distroless/static-debian13

WORKDIR /
COPY --from=build-base /app/bin/blog-release /blog

CMD ["./blog"]
