#!/bin/sh

if [ "$1" = "api" ]; then
    echo "Starting API server..."
    exec npm start
elif [ "$1" = "worker" ]; then
    echo "Starting worker..."
    exec npm run start:worker
else
    echo "Unknown command: $1"
    echo "Use 'api' or 'worker'"
    exit 1
fi