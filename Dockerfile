
FROM gcr.io/cloud-builders/go as build

ENV GOPATH /go
ENV GO111MODULE on

WORKDIR ${GOPATH}/src

# Copy sources
COPY . .

RUN wget -O /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.3.1/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -installsuffix "static" ./examples/server/main.go

FROM gcr.io/distroless/static:nonroot

COPY --from=build /go/bin/main /bin/greeter
COPY --from=build /bin/grpc_health_probe /bin/grpc_health_probe

ENTRYPOINT ["/bin/greeter"]
