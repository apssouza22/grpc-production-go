# GRPC for production [![Build Status](https://travis-ci.org/apssouza22/grpc-server-go.svg?branch=master)](https://travis-ci.org/apssouza22/grpc-server-go) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=apssouza22_grpc-server-go&metric=alert_status)](https://sonarcloud.io/dashboard?id=apssouza22_grpc-server-go) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=apssouza22_grpc-server-go&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=apssouza22_grpc-server-go)
This GRPC project provides:

- Health check service
- Shutdown hook
- Keep alive params
- Server interceptors(auth, request canceled, execution time)
- Idempotent check
- Server example
- Client example (timeout, healthcheck)
- Client interceptor (timeout log)
- Unit tests using a In memory communication between client and server
