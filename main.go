package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/cloudfoundry/noaa/consumer"

	"github.com/cf-platform-eng/firehose-nozzle/auth"
	"github.com/cf-platform-eng/firehose-nozzle/config"
	"github.com/cf-platform-eng/firehose-nozzle/nozzle"
	"github.com/cf-platform-eng/firehose-nozzle/writernozzle"
)

func main() {
	cf_lager.AddFlags(flag.CommandLine)
	flag.Parse()

	logger, _ := cf_lager.New("firehose-logger")
	logger.Info("Running firehose-nozzle")

	config, err := config.Parse()
	if err != nil {
		logger.Fatal("Unable to build config from environment", err)
	}

	fetcher := auth.NewUAATokenFetcher(config.UAAURL, config.Username, config.Password, true)
	token, err := fetcher.FetchAuthToken()
	if err != nil {
		logger.Fatal("Unable to fetch token", err)
	}

	consumer := consumer.New(config.TrafficControllerURL, &tls.Config{
		InsecureSkipVerify: config.InsecureSkipVerify,
	}, nil)
	events, errors := consumer.Firehose(config.FirehoseSubscriptionID, token)

	writerEventSerializer := writernozzle.NewWriterEventSerializer()
	writerClient := writernozzle.NewWriterClient(os.Stderr)
	logger.Info(fmt.Sprintf("Forwarding events: %s", config.SelectedEvents))
	forwarder := nozzle.NewForwarder(
		writerClient, writerEventSerializer,
		config.SelectedEvents, events, errors, logger,
	)
	err = forwarder.Run(time.Second)
	if err != nil {
		logger.Fatal("Error forwarding", err)
	}
}
