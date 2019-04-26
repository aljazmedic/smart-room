#!/usr/bin/bash
docker-compose up -d
echo "Running in background"
docker-compose ps | cat