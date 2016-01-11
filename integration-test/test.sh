#!/bin/bash
docker-compose stop && docker-compose rm -f && docker-compose up -d
authstatus=0
while [ $authstatus -ne 401 ]
do
  sleep 1
  authstatus=$(curl -s -o /dev/null -w "%{http_code}" http://$DH/authentication-service/oauth/token)
  echo "got status $authstatus from authentication-service"
done
echo "authentication-service ready!"
go test
docker-compose stop && docker-compose rm -f
