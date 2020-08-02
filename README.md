# GRPC for production [![Build Status](https://travis-ci.org/apssouza22/grpc-server-go.svg?branch=master)](https://travis-ci.org/apssouza22/grpc-server-go) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=apssouza22_grpc-server-go&metric=alert_status)](https://sonarcloud.io/dashboard?id=apssouza22_grpc-server-go) [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=apssouza22_grpc-server-go&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=apssouza22_grpc-server-go)

Read more about the project [here!](https://medium.com/@alexsandrosouza/grpc-for-production-go-2f62f334824)

This project abstracts away the details of the GRPC server and client configuration. 

Here are the main features:
- Health check service — We use the grpc_health_probe utility which allows you to query health of gRPC services that expose service their status through the gRPC Health Checking Protocol.
- Shutdown hook — The library registers a shutdown hook with the GRPC server to ensure that the application is closed gracefully on exit
- Keep alive params — Keepalives are an optional feature but it can be handy to signal how the persistence of the open connection should be kept for further messages
- In memory communication between client and server, helpful to write unit and integration tests. When writing integration tests we should avoid having the networking element from your test as it is slow to assign and release ports.
- Server and client builder for uniform object creation
- Added ability to recover the system from a service panic
- Added ability to add multiple interceptors in order
- Added client tracing metadata propagation
- Handy Server interceptors(Authentication, request cancelled, execution time, panic recovery)
- Handy Client interceptors(Timeout logs, Tracing, propagate headers)
- Secure connection with self signed certificate
 
 ## Examples
 
 Please refer to the /examples folder
 
 
