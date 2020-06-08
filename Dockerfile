FROM golang:1.14

WORKDIR /go/src/go_api
COPY . .

RUN go get -u -v github.com/swaggo/swag/cmd/swag
RUN swag init
RUN go get -v

CMD ["go_api"]
