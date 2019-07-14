package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var out *logrus.Logger
var client *redis.Client

type value struct {
	Key   string
	Value string
	Avail []details
}

type details struct {
	Start    string
	Duration int
}

func isAllowed(ip string, list []string) bool {
	for _, l := range list {
		out.Info(l)
		if l == ip {
			return true
		}
	}
	return false
}

func get(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)

	if isAllowed(strings.Split(r.RemoteAddr, ":")[0], strings.Split(os.Getenv("ALLOWED"), ":")) {
		out.Warn("failed request ", r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	out.Info("new request ", r.RemoteAddr)

	val, err := client.Get(p["key"]).Result()
	if err == redis.Nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		out.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("unable get value from cache")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, val)
}

func output(dev bool) {
	out = logrus.New()
	out.SetFormatter(&logrus.JSONFormatter{})
	out.Out = os.Stdout

	if !dev {
		file, err := os.OpenFile("output/get", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			out.Out = file
		}
	}
}

func main() {
	output(true)

	r := mux.NewRouter()
	r.HandleFunc("/{key}", get).Methods("GET")

	out.Info("creating cache client")
	client = redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})

	out.Info("ping cache")
	_, err := client.Ping().Result()
	if err != nil {
		out.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("unable to ping cache")
	}

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		out.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("unable to start service")
	}
}
