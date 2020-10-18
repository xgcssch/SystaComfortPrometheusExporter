FROM golang:1.15

WORKDIR /go/src/app
COPY . .

RUN go build -v ./... 
RUN go install -v ./...

CMD ["app"]