#!/bin/bash

# Trap SIGINT (Ctrl+C) to kill all background processes
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

echo "Starting Echelon Services..."

# Start Delta (API)
echo "[Delta] Starting on port 8080..."
go run pkg/delta/main.go &
DELTA_PID=$!

# Start Bonfire (Websocket)
echo "[Bonfire] Starting on port 8081..."
go run pkg/bonfire/main.go &
BONFIRE_PID=$!

# Wait for both processes
wait $DELTA_PID $BONFIRE_PID
