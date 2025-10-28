# --- Build stage for Go API ---
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o sawapi main.go

# --- SAW + Cryptol base layers ---
FROM ghcr.io/galoisinc/saw:nightly AS sawsrc
FROM ghcr.io/galoisinc/cryptol:3.2.0 AS cryptolsrc

# --- Main runtime image ---
FROM ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive

# Install Clang/LLVM and basic tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates wget curl gnupg software-properties-common build-essential \
    z3 git vim && rm -rf /var/lib/apt/lists/*

# Copy SAW & Cryptol binaries
COPY --from=sawsrc /usr/local /usr/local
COPY --from=cryptolsrc /usr/local /usr/local
ENV PATH="/usr/local/bin:${PATH}"

# Copy Go binary
COPY --from=builder /app/sawapi /usr/local/bin/sawapi

# Copy example workspace
WORKDIR /work
COPY example ./example

EXPOSE 8443
CMD ["/usr/local/bin/sawapi"]