#!/bin/sh
cd integration-test
docker-compose stop && docker-compose rm -f && docker-compose up -d
sleep 60 && go test
cd ..

