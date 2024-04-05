#!/usr/bin/env bash
set -e

CMD="./sidekick --port=8080 --app-port=8888 $* \
    php -S localhost:8888 -t example"

echo "$CMD"
$CMD
