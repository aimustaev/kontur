#!/bin/sh

if [ "$1" = "api" ]; then
    echo "Starting API server..."
    exec node src/index.js
elif [ "$1" = "worker" ]; then
    echo "Starting worker..."
    exec node ./src/worker.js
else
    echo "Unknown command: $1"
    echo "Use 'api' or 'worker'"
    exit 1
fi