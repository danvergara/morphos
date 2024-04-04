# Build the application from source
FROM golang:1.21 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN apt-get update \
   && apt-get install -y --no-install-recommends fonts-recommended libvips-dev \
   && apt-get autoremove -y \
   && apt-get purge -y --auto-remove \
   && rm -rf /var/lib/apt/lists/*

COPY go.* ./
RUN go mod download

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o morphos .

# Deploy the application binary into a lean image
FROM debian:trixie-slim AS release

WORKDIR /

RUN apt-get update \
   && apt-get install -y --no-install-recommends default-jre libreoffice libreoffice-java-common libvips42 \
   && apt-get autoremove -y \
   && apt-get purge -y --auto-remove \
   && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/morphos /bin/morphos
COPY --from=builder /usr/share/fonts /usr/share/fonts

ENV FONTCONFIG_PATH /usr/share/fonts
# memory arena allocation for libvips - see https://github.com/davidbyttow/govips README
ENV MALLOC_ARENA_MAX 2

EXPOSE 8080

ENTRYPOINT ["/bin/morphos"]
