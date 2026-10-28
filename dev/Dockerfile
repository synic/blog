FROM docker.io/library/golang:1.27.3-trixie

RUN apt-get update \
    && apt-get install -y --no-install-recommends webp nodejs npm python3 python3-venv \
    && rm -rf /var/lib/apt/lists/*
