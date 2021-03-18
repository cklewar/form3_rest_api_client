#!/bin/sh
until $(curl --output /dev/null --silent --location --request GET --fail 'http://accountapi:8080/v1/health' --header 'Content-Type: application/vnd.api+json');do printf '.';sleep 5;done
go test -v ./api/client/ -coverprofile cover.out
go run ./examples/main.go