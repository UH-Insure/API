# API
go api to run cryptol and saw files from collab



# SAW + Cryptol Remote API

A Go HTTPS API for running SAW and Cryptol commands from remote clients
(e.g. Google Colab) via Cloudflare tunnel.

## Features
- Run `.cry` and `.saw` files remotely
- HTTPS with Cloudflare tunnel
- Dockerized with full SAW + Cryptol toolchain

## Run locally
```bash
docker build -t saw-cryptol-api .
docker run -p 8443:8443 saw-cryptol-api
```