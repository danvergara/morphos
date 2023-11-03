# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH}  go build -o morphos .

# Deploy the application binary into a lean image
FROM debian:bookworm-slim AS release

WORKDIR /

COPY --from=builder /app/morphos /bin/morphos

EXPOSE 8080

ENTRYPOINT ["/bin/morphos"]
