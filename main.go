package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/cloudfoundry/noaa/consumer"

	"github.com/cf-platform-eng/firehose-nozzle/auth"
	"github.com/cf-platform-eng/firehose-nozzle/config"
	"github.com/cf-platform-eng/firehose-nozzle/nozzle"
	"github.com/cf-platform-eng/firehose-nozzle/writernozzle"
)

func main() {
	logger := log.New(os.Stdout, ">>> ", 0)

	conf, err := config.Parse()
	if err != nil {
		logger.Fatal("Unable to build config from environment", err)
	}

	fetcher := auth.NewAPITokenFetcher(conf.APIURL, conf.Username, conf.Password, true)
	token, err := fetcher.FetchAuthToken()
	if err != nil {
		logger.Fatal("Unable to fetch token", err)
	}

	noaaConsumer := consumer.New(conf.TrafficControllerURL, &tls.Config{
		InsecureSkipVerify: conf.SkipSSL,
	}, nil)
	events, errors := noaaConsumer.Firehose(conf.FirehoseSubscriptionID, token)

	writerEventSerializer := writernozzle.NewWriterEventSerializer()
	writerClient := writernozzle.NewWriterClient(os.Stdout)
	logger.Printf("Forwarding events: %s", conf.SelectedEvents)
	forwarder := nozzle.NewForwarder(
		writerClient, writerEventSerializer,
		conf.SelectedEvents, events, errors, logger,
	)
	err = forwarder.Run(time.Second)
	if err != nil {
		logger.Fatal("Error forwarding", err)
	}
}
