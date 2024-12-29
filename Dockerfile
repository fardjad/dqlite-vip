# syntax=docker/dockerfile:1

FROM golang:1.23-bookworm AS builder

RUN echo "deb http://deb.debian.org/debian trixie main" > /etc/apt/sources.list.d/trixie.list \
  && echo "Package: *\nPin: release a=trixie\nPin-Priority: 100" > /etc/apt/preferences.d/trixie \
  && echo "Package: just\nPin: release a=trixie\nPin-Priority: 500" >> /etc/apt/preferences.d/trixie \
  && apt-get update \
  && apt-get install -y --no-install-recommends just \
  && rm -rf /var/lib/apt/lists/* \
  && rm /etc/apt/sources.list.d/trixie.list \
  && rm /etc/apt/preferences.d/trixie

RUN apt-get update && apt-get install -y --no-install-recommends \
  build-essential \
  automake \
  libtool \
  gettext \
  autopoint \
  tclsh \
  tcl \
  libsqlite3-dev \
  pkg-config \
  git \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY hack/ hack/
RUN ./hack/static-dqlite.sh
COPY . .
RUN just build-static

FROM scratch

COPY --from=builder --chmod=755 /app/bin/static/dqlite-vip /dqlite-vip

ENTRYPOINT ["/dqlite-vip"]