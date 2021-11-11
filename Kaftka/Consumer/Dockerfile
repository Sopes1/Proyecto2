FROM golang:latest

RUN mkdir /consumerKaftka



WORKDIR /consumerKaftka
#Para exportar modulos de go
RUN export GO111MODULE=on
RUN go get go.mongodb.org/mongo-driver
RUN go get -u gopkg.in/confluentinc/confluent-kafka-go.v1/kafka
RUN go get github.com/go-redis/redis/v8


#Compilar go

COPY . . 

RUN go build

ENTRYPOINT [ "/consumerKaftka/quickstart" ]