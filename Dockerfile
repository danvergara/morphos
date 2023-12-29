FROM ubuntu:22.04 AS base

ARG GO_VERSION
ENV GO_VERSION=${GO_VERSION}

RUN apt-get update
RUN apt-get install -y wget git gcc

RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"

RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR $GOPATH

FROM base AS builder

RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:strukturag/libde265 && \
    add-apt-repository ppa:strukturag/libheif && \
    apt-get install -y --no-install-recommends cmake \
      make \
      pkg-config \
      x265 \
      libx265-dev \
      libde265-dev \
      libjpeg-dev \
      libtool \
      zlib1g-dev \
      libaom-dev \
      libheif1 \
      libheif-dev && \
      apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags='-s -w' -trimpath -o /app/morphos .

ENTRYPOINT ["/app/morphos"]
