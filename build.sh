#!/bin/bash
set -ev
go test
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o out/main .
docker build -t openservice/lecture-service:latest .
cd integration-test
docker-compose stop && docker-compose rm -f && docker-compose up -d
sleep 30 && go test
cd ..
if [ "${TRAVIS_PULL_REQUEST}" = "false" ] && [ "${TRAVIS_REPO_SLUG}" = "InteractiveLecture/lecture-service" ] ; then
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD" -e="$DOCKER_EMAIL"
  docker push openservice/lecture-service:latest
fi
