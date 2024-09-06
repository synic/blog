FROM golang:1.23.0-alpine3.20 AS build-base

WORKDIR /app

COPY . .

ENV GOPATH=/go
ENV PATH="${PATH}:/go/bin"
RUN set -x \
    && apk add --no-cache nodejs=20.15.1-r0 npm=10.8.0-r0 make=4.4.1-r2  \
    && make release \
    && rm -rf node_modules

FROM gcr.io/distroless/static-debian12:9efbcaacd8eac4960b315c502adffdbf3398ce62

WORKDIR /
COPY --from=build-base /app/bin/blog-release /blog
COPY --from=build-base /app/assets /assets
COPY --from=build-base /app/articles /articles

CMD ["./blog", "serve"]
