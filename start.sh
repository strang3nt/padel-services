#!/bin/bash

if [ -f .env ]; then
    set -a
    source .env
    set +a
    echo "✅ Environment variables loaded from .env"
else
    echo "⚠️  .env file not found, using system defaults"
fi

echo "🚀 Starting Go application..."
go run ./cmd/tgbot
