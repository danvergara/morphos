FROM ubuntu:22.04 AS build

ARG GO_VERSION TARGETOS TARGETARCH

ENV GO_VERSION=${GO_VERSION}
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:strukturag/libde265 && \
    add-apt-repository ppa:strukturag/libheif && \
    apt-get install -y --no-install-recommends cmake \
      wget \
      git \
      gcc \
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
      apt-get autoremove -y && \
      apt-get purge -y --auto-remove && \
    rm -rf /var/lib/apt/lists/*

RUN wget --no-check-certificate -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" && \
    tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz" && \
    rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags='-s -w' -trimpath -o /app/morphos .

FROM ubuntu:22.04

WORKDIR /app

COPY --from=build /app/morphos .

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
      apt-get autoremove -y && \
      apt-get purge -y --auto-remove && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 8080

ENTRYPOINT ["/app/morphos"]
