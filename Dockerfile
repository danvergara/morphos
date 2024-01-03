# Build the application from source
FROM golang:1.21 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o morphos .

# Deploy the application binary into a lean image
FROM debian:bookworm-slim AS release

WORKDIR /

COPY --from=builder /app/morphos /bin/morphos

EXPOSE 8080

ENTRYPOINT ["/bin/morphos"]
