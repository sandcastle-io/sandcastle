FROM debian:bookworm-slim AS nsjail-builder

RUN apt-get update && apt-get install -y \
    build-essential \
    bison \
    flex \
    pkg-config \
    libnl-route-3-dev \
    libprotobuf-dev \
    protobuf-compiler \
    git \
    ca-certificates

RUN git clone https://github.com/google/nsjail.git /nsjail && \
    cd /nsjail && \
    make

FROM golang:1.25-bookworm AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /kube-sandcastle-worker ./cmd/sandcastle-worker

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    libprotobuf32 \
    libnl-route-3-200 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=nsjail-builder /nsjail/nsjail /usr/bin/nsjail

COPY --from=go-builder /kube-sandcastle-worker /app/worker

RUN groupadd -g 2000 sandcastle && \
    useradd -u 2000 -g sandcastle -m sandcastle

WORKDIR /app

CMD ["/app/worker"]
