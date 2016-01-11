#!/bin/bash
set -ev
go test
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/main .
docker build -t openservice/lecture-service:latest .
