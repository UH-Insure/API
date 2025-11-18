# ============================================
# SAW + Cryptol + Solvers + Go API + SQLite
# ============================================

# Pull upstream SAW & Cryptol distributions
FROM ghcr.io/galoisinc/saw:nightly AS sawsrc
FROM ghcr.io/galoisinc/cryptol:3.2.0 AS cryptolsrc

# Go builder
FROM golang:1.22 AS gobuild
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o sawapi main.go

# Main runtime container
FROM ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive

# Install base tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates wget curl gnupg lsb-release software-properties-common \
    build-essential clang-16 clang-tools-16 z3 sqlite3 git python3 python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Install C compiler support
RUN update-alternatives --install /usr/bin/cc  cc  /usr/bin/clang-16  100 && \
    update-alternatives --install /usr/bin/c++ c++ /usr/bin/clang++-16 100

# Copy SAW & Cryptol toolchains
COPY --from=sawsrc /usr/local /usr/local
COPY --from=cryptolsrc /usr/local /usr/local

ENV PATH="/usr/local/bin:${PATH}"

# Create runtime workspace
RUN mkdir -p /work && mkdir -p /data

# Copy Go API server
COPY --from=gobuild /app/sawapi /usr/local/bin/sawapi

# Copy example files for testing
COPY example /work/example

# SQLite DB
COPY internal/schema.sql /data/schema.sql

# Initialize DB at container start if missing
COPY internal/init-db.sh /usr/local/bin/init-db.sh
RUN chmod +x /usr/local/bin/init-db.sh

EXPOSE 8443
CMD ["/usr/local/bin/init-db.sh"]