package main

import (
	"encoding/gob"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var out *logrus.Logger

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
		if l == ip {
			return true
		}
	}
	return false
}

func set(w http.ResponseWriter, r *http.Request) {
	if isAllowed(strings.Split(r.RemoteAddr, ":")[0], strings.Split(os.Getenv("ALLOWED"), ":")) {
		out.Warn("failed request ", r.RemoteAddr)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	out.Info("new request ", r.RemoteAddr)

	var val value
	w.Header().Set("content-type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&val)
	if err != nil {
		out.WithFields(logrus.Fields{
			"error": err,
		}).Warn("unable to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, err := os.Create("data/" + val.Key)
	if err != nil {
		out.WithFields(logrus.Fields{
			"error": err,
		}).Warn("unable to save gob file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encoder := gob.NewEncoder(file)
	encoder.Encode(val)

	file.Close()

	out.WithFields(logrus.Fields{
		"data": val,
	}).Info("successfully updated")

}

func output(dev bool) {
	out = logrus.New()
	out.SetFormatter(&logrus.JSONFormatter{})
	out.Out = os.Stdout

	if !dev {
		file, err := os.OpenFile("output/put", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			out.Out = file
		}
	}
}

func main() {
	output(true)

	r := mux.NewRouter()
	r.HandleFunc("/", set).Methods("POST")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		out.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("unable to start service")
	}
}
