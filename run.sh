#!/usr/bin/env bash
set -e
./sidekick \
    --port=8080 \
    --app-port=8888 \
    php -S localhost:8888 -t example
