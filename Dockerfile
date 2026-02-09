FROM golang:1.26rc3-trixie AS build-base

WORKDIR /app
COPY . .

RUN go run github.com/magefile/mage@2385abb build:release
FROM gcr.io/distroless/static-debian13:972618ca78034aaddc55864342014a96b85108c607372f7cbd0dbd1361f1d841

WORKDIR /
COPY --from=build-base /app/bin/blog-release /blog

CMD ["./blog"]
