#!/bin/bash

echo "Buildig game_server app"

CGO_ENABLED=0 GOOS=linux go build -o game_server .

echo "Starting game_server"

./game_server