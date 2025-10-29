# Dockerfile
FROM ghcr.io/galoisinc/cryptol-remote-api:latest

# Optional: health check (useful for container orchestration)
HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/info || exit 1