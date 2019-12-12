# GRPC for production [![Build Status](https://travis-ci.org/apssouza22/grpc-server-go.svg?branch=master)](https://travis-ci.org/apssouza22/grpc-server-go) [![Code Climate](https://codeclimate.com/github/apssouza22/grpc-server-go.png)](https://codeclimate.com/github/apssouza22/grpc-server-go) [![Test Coverage](https://api.codeclimate.com/v1/badges/d999a5f1311bd806b345/test_coverage)](https://codeclimate.com/github/apssouza22/grpc-server-go/test_coverage) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=apssouza22_grpc-server-go&metric=alert_status)](https://sonarcloud.io/dashboard?id=apssouza22_grpc-server-go)
This GRPC server base provides:

- Health check service
- Shutdown hook
- Keep alive params
- Server interceptors(auth, request canceled, execution time)
- Idempotent check
- Server example
- Client example (timeout, healthcheck)
- Client interceptor (timeout log)
- Unit tests using a In memory communication between client and server
