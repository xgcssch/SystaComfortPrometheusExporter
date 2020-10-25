# Build with
# docker build -t SystaComfortPrometheusExporter -f build/docker/Dockerfile .
# from root directory
FROM golang:alpine as builder

WORKDIR /go/app
COPY . .

RUN mkdir bin && \
    GOOS=linux go build -v -o bin ./... && \
    go install -v ./...

FROM alpine:latest
# docker run -it systa /bin/ash

EXPOSE 22460/udp
EXPOSE 2112/tcp

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/app/bin/SystaComfortPrometheusExporter .

CMD ["./SystaComfortPrometheusExporter"]
