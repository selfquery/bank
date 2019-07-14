package main

import (
	"encoding/gob"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/jasonlvhit/gocron"
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

func process(val value, cache *redis.Client) {
	ct := time.Now().Format("3:04pm")
	for t := range val.Avail {
		if ct == val.Avail[t].Start {
			err := cache.Set(val.Key, val.Value, time.Duration(val.Avail[t].Duration)*time.Second).Err()
			if err != nil {
				out.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("unable to move key value to cache")
			}
			out.WithFields(logrus.Fields{
				"key": val,
			}).Info("moved key to cache")
		}
	}
}

func walk(cache *redis.Client) {
	out.Info("walking data directory")

	files, err := ioutil.ReadDir("./data")
	if err != nil {
		out.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("unable to read directory")
	}

	for _, f := range files {
		var val value

		file, err := os.Open("./data/" + f.Name())
		if err != nil {
			out.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("unable to open key file")
		}

		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&val)
		if err != nil {
			out.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("unable to decode key file")
		}

		file.Close()

		process(val, cache)
	}
}

func output(dev bool) {
	out = logrus.New()
	out.SetFormatter(&logrus.JSONFormatter{})
	out.Out = os.Stdout

	if !dev {
		file, err := os.OpenFile("output/process", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			out.Out = file
		}
	}
}

func main() {
	output(true)

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

	gocron.Every(1).Minute().Do(walk, client)
	<-gocron.Start()
}
