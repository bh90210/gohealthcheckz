FROM golang:1.16
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "testapp" ./examples/http/main.go

FROM alpine:3.13.5
COPY --from=0 /app/testapp /bin/
# RUN apk -Uuv add ca-certificates