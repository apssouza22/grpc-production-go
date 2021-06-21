# GRPC for production [![Build Status](https://travis-ci.org/apssouza22/grpc-server-go.svg?branch=master)](https://travis-ci.org/apssouza22/grpc-server-go) [! [![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=apssouza22_grpc-server-go&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=apssouza22_grpc-server-go)

Read more about the project [here!](https://dev.to/apssouza22/grpc-for-production-golang-1611)

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
- Client TLS with insecure connection support 


---
## Free Advanced Java Course
I am the author of the [Advanced Java for adults course](https://www.udemy.com/course/advanced-java-for-adults/?referralCode=8014CCF0A5A931ADED5F). This course contains advanced and not conventional lessons. In this course, you will learn to think differently from those who have a limited view of software development. I will provoke you to reflect on decisions that you take in your day to day job, which might not be the best ones. This course is for middle to senior developers and we will not teach Java language features but how to lead complex Java projects. 

This course's lectures are based on a Trading system, an opensource project hosted on my [Github](https://github.com/apssouza22/trading-system).

---

 
 ## Examples
 
 Please refer to the /examples folder
 
 
