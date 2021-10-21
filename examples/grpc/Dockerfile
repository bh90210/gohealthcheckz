FROM golang:1.16
COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0 
ENV GOOS=linux 
ENV GOARCH=amd64 
RUN go mod init testapp
RUN go mod tidy
RUN go build -o "testapp" ./main.go

FROM alpine:3.13.5

COPY --from=0 /app/testapp /bin/

RUN apk add wget
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 
RUN chmod +x /bin/grpc_health_probe  