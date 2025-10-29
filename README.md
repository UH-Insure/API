# API
all this does is connect the cryptol api to collab. The go files connected the compiler on my computer to collab for saw and cryptol, but this was deemed unecessary so now all it does is connect cryptol api to collab.



# old stuff

## SAW + Cryptol Remote API

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


## Server setup

git pull origin main
sudo podman build -t localhost/saw-cryptol-api:latest .
sudo podman run -d --replace \
  --name saw-cryptol-api \
  -p 8443:8443 \
  -v $(pwd)/work:/work:Z \
  localhost/saw-cryptol-api:latest