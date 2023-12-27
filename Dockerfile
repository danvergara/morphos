FROM golang:1.21-alpine AS builder

RUN apk add --no-cache build-base \
  pkgconf \
  libgcc \
  libstdc++ \
  libwebp-dev \
  libsharpyuv \
  x265-libs \
  libde265 \
  libde265-dev \
  musl \
  aom-libs \
  libheif \
  libheif-dev

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o morphos .

EXPOSE 8080
ENTRYPOINT ["/app/morphos"]
