FROM golang:latest

ADD . /audit
WORKDIR /audit

EXPOSE 8080

RUN go get "github.com/go-redis/redis"
RUN go get "github.com/sirupsen/logrus"
RUN go get "github.com/jasonlvhit/gocron"

RUN go build -o bin .
CMD ["./bin"]