FROM golang:latest

RUN mkdir /grpcClient

WORKDIR /grpcClient

RUN export GO111MODULE=on
COPY . .
RUN go mod download
#RUN go mod init github.com/Sopes1
#RUN go get -u github.com/gorilla/mux
#RUN go get google.golang.org/grpc

EXPOSE 3200


CMD ["go","run","client.go"]

# docker build -t grpcclient .

## docker run -it --rm --name grpcclienttt -p 3200:3200 grpcclient:latest