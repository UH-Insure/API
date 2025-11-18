# ============================================
#  SAW + Cryptol + Full Solver Suite + Go API
# ============================================

##### Upstream SAW & Cryptol
FROM ghcr.io/galoisinc/saw:nightly AS sawsrc
FROM ghcr.io/galoisinc/cryptol:3.2.0 AS cryptolsrc

##### Go Builder
FROM golang:1.22 AS gobuild
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o sawapi main.go


##############################################
#       MAIN RUNTIME IMAGE
##############################################
FROM ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive

# ------------------------------
# Install Clang/LLVM + utilities
# ------------------------------
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates wget curl gnupg lsb-release software-properties-common \
  && curl -fsSL https://apt.llvm.org/llvm-snapshot.gpg.key \
        | gpg --dearmor -o /usr/share/keyrings/llvm.gpg \
  && echo "deb [signed-by=/usr/share/keyrings/llvm.gpg] http://apt.llvm.org/jammy/ llvm-toolchain-jammy-16 main" \
        > /etc/apt/sources.list.d/llvm16.list \
  && apt-get update \
  && apt-get install -y --no-install-recommends \
       clang-16 clang-tools-16 clang-format-16 lldb-16 lld-16 \
       llvm-16 libc++-16-dev libc++abi-16-dev \
       build-essential cmake git zsh unzip python3 python3-pip sqlite3 \
       z3 bash \
  && rm -rf /var/lib/apt/lists/*

# Make clang-16 default
RUN update-alternatives --install /usr/bin/cc  cc  /usr/bin/clang-16  100 && \
    update-alternatives --install /usr/bin/c++ c++ /usr/bin/clang++-16 100


# ------------------------------
# Install SMT Solvers
# ------------------------------
RUN apt-get update && apt-get install -y --no-install-recommends \
      libgmp-dev libffi-dev libreadline-dev \
  && mkdir -p /opt/solvers && cd /opt/solvers \
  \
  # ---- ABC ----
  && git clone https://github.com/berkeley-abc/abc.git \
  && cd abc && make -j$(nproc) && cp abc /usr/local/bin/ && cd .. \
  \
  # ---- Boolector ----
  && git clone https://github.com/Boolector/boolector.git \
  && cd boolector && ./contrib/setup-lingeling.sh && ./contrib/setup-btor2tools.sh \
  && ./configure.sh && cd build && make -j$(nproc) install && cd ../.. \
  \
  # ---- CVC4 ----
  && wget --no-verbose https://cvc4.cs.stanford.edu/downloads/builds/x86_64-linux-opt/cvc4-1.8-x86_64-linux-opt \
  && chmod +x cvc4-1.8-x86_64-linux-opt && mv cvc4-1.8-x86_64-linux-opt /usr/local/bin/cvc4 \
  \
  # ---- CVC5 ----
  && wget --no-verbose https://github.com/cvc5/cvc5/releases/download/cvc5-1.3.1/cvc5-Linux-x86_64-static.zip \
  && unzip cvc5-Linux-x86_64-static.zip \
  && cp cvc5-Linux-x86_64-static/bin/cvc5 /usr/local/bin/cvc5 \
  && rm -rf cvc5-Linux-x86_64-static* \
  \
  # ---- MathSAT ----
  && wget --no-verbose https://mathsat.fbk.eu/release/mathsat-5.6.12-linux-x86_64.tar.gz -O mathsat.tar.gz \
  && tar -xzf mathsat.tar.gz \
  && cp mathsat-*/bin/mathsat /usr/local/bin/mathsat \
  && rm -rf mathsat* \
  \
  # ---- Yices ----
  && wget --no-verbose https://yices.csl.sri.com/releases/2.6.4/yices-2.6.4-x86_64-pc-linux-gnu.tar.gz -O yices.tar.gz \
  && tar -xzf yices.tar.gz \
  && cp yices-*/bin/yices /usr/local/bin/yices \
  && rm -rf yices* \
  \
  && rm -rf /var/lib/apt/lists/*


# ------------------------------
# Copy SAW & Cryptol Binaries
# ------------------------------
COPY --from=sawsrc /usr/local /usr/local
COPY --from=cryptolsrc /usr/local /usr/local
ENV PATH="/usr/local/bin:${PATH}"


# ------------------------------
# Workspace & Data Setup
# ------------------------------
RUN mkdir -p /work /work/files /data

COPY example /work/example
COPY internal/storage/schema.sql /data/schema.sql
COPY internal/storage/init-db.sh /usr/local/bin/init-db.sh
RUN chmod +x /usr/local/bin/init-db.sh


# ------------------------------
# Copy Go API
# ------------------------------
COPY --from=gobuild /app/sawapi /usr/local/bin/sawapi


EXPOSE 8443

# Start DB then API
CMD ["/usr/local/bin/init-db.sh"]