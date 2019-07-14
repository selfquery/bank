FROM golang:latest

ADD . /get
WORKDIR /get

EXPOSE 8080

RUN go get "github.com/gorilla/mux"
RUN go get "github.com/sirupsen/logrus"
RUN go get "github.com/go-redis/redis"

RUN go build -o bin ./src
CMD ["./bin"]