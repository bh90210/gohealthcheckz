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
