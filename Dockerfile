FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get ./...

CMD go run main.go