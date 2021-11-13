FROM golang:latest

RUN mkdir /grpcServer

WORKDIR /grpcServer

RUN export GO111MODULE=on
COPY . .
RUN go mod download
#RUN go mod init github.com/Sopes1
#RUN go get github.com/streadway/amqp
#RUN go get google.golang.org/grpc


CMD ["go","run","server.go"]

## docker build -t grpcserver .

## docker run -it --rm --name grpcserver grpcserver:latest