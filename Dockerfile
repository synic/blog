FROM golang:1.26.2-trixie AS build-base

WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y webp && rm -rf /var/lib/apt/lists/*
RUN go tool github.com/magefile/mage build:release

FROM gcr.io/distroless/static-debian13

WORKDIR /
COPY --from=build-base /app/bin/blog-release /blog
COPY --from=build-base /app/static /static
COPY --from=build-base /app/migrations /migrations
ENV STATIC_DIR=/static
ENV MIGRATIONS_DIR=/migrations

CMD ["./blog"]
