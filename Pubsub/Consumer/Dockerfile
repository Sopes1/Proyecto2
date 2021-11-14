FROM golang:latest

RUN mkdir /consumerPubsub



WORKDIR /consumerPubsub
#Para exportar modulos de go
RUN export GO111MODULE=on
RUN go get cloud.google.com/go/pubsub
RUN go get github.com/go-redis/redis/v8
RUN go get go.mongodb.org/mongo-driver

#Compilar go

COPY . . 
ENV GOOGLE_APPLICATION_CREDENTIALS="./key.json"
RUN go build

ENTRYPOINT [ "/consumerPubsub/quickstart" ]