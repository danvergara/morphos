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

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags='-s -w' -trimpath -o /dist/morphos .
RUN ldd /dist/morphos | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /dist/%

FROM scratch
COPY --from=builder /dist /
USER 65534

EXPOSE 8080
ENTRYPOINT ["/morphos"]
