FROM golang:latest

ADD . /put
WORKDIR /put

EXPOSE 8080

RUN go get "github.com/gorilla/mux"
RUN go get "github.com/sirupsen/logrus"

RUN go build -o bin .
CMD ["./bin"]