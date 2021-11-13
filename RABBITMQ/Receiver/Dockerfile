FROM golang:latest

RUN mkdir /rabbit

WORKDIR /rabbit

RUN export GO111MODULE=on

RUN go get github.com/rabbitmq/amqp091-go
RUN go get go.mongodb.org/mongo-driver
RUN go get go.mongodb.org/mongo-driver/mongo/options
RUN go get github.com/go-redis/redis/v8

COPY . .

RUN go build

ENTRYPOINT [ "/rabbit/Proyecto2" ]

# docker build -t workerrabbit .
# docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.9-management
# docker run -it --rm --env QRHOST=localhost --env MHOST=10.128.0.5 --env RHOST=10.128.0.8 --name rabbitworker 6dc82dcd6e6f