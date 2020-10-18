FROM golang:alpine as builder

WORKDIR /go/app
COPY . .

RUN mkdir bin && \
    CGO_ENABLED=0 GOOS=linux go build -v -o bin ./... && \
    go install -v ./...

FROM alpine:latest
# docker run -it systa /bin/ash

EXPOSE 22460/udp

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/app/bin/SystaComfortPrometheusExporter .

CMD ["./SystaComfortPrometheusExporter"]
