package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"
	"errors"

	"github.com/cloudfoundry/noaa/consumer"

	"github.com/cf-platform-eng/firehose-nozzle/api"
	"github.com/cf-platform-eng/firehose-nozzle/config"
	"github.com/cf-platform-eng/firehose-nozzle/nozzle"
	"github.com/cf-platform-eng/firehose-nozzle/uaa"
	"github.com/cf-platform-eng/firehose-nozzle/writernozzle"
)

func main() {
	logger := log.New(os.Stdout, ">>> ", 0)

	conf, err := config.Parse()
	if err != nil {
		logger.Fatal("Unable to build config from environment", err)
	}

	var token, trafficControllerURL string
	if conf.APIURL != "" {
		logger.Printf("Fetching auth token via API: %v\n", conf.APIURL)

		fetcher, err := api.NewAPIClient(conf.APIURL, conf.Username, conf.Password, conf.SkipSSL)
		if err != nil {
			logger.Fatal("Unable to build API client", err)
		}
		token, err = fetcher.FetchAuthToken()
		if err != nil {
			logger.Fatal("Unable to fetch token via API", err)
		}

		trafficControllerURL = fetcher.FetchTrafficControllerURL()
		if trafficControllerURL == "" {
			logger.Fatal("trafficControllerURL from client was blank")
		}
		if conf.ResolveAppMetaData {
			//do something clever
		}
	} else if conf.UAAURL != "" {
		logger.Printf("Fetching auth token via UAA: %v\n", conf.UAAURL)

		trafficControllerURL = conf.TrafficControllerURL
		if trafficControllerURL == "" {
			logger.Fatal(errors.New("NOZZLE_TRAFFIC_CONTROLLER_URL is required when authenticating via UAA"))
		}

		fetcher := uaa.NewUAATokenFetcher(conf.UAAURL, conf.Username, conf.Password, conf.SkipSSL)
		token, err = fetcher.FetchAuthToken()
		if err != nil {
			logger.Fatal("Unable to fetch token via UAA", err)
		}
	} else {
		logger.Fatal(errors.New("One of NOZZLE_API_URL or NOZZLE_UAA_URL are required"))
	}

	logger.Printf("Consuming firehose: %v\n", trafficControllerURL)
	noaaConsumer := consumer.New(trafficControllerURL, &tls.Config{
		InsecureSkipVerify: conf.SkipSSL,
	}, nil)
	events, errs := noaaConsumer.Firehose(conf.FirehoseSubscriptionID, token)

	writerEventSerializer := writernozzle.NewWriterEventSerializer()
	writerClient := writernozzle.NewWriterClient(os.Stdout)
	logger.Printf("Forwarding events: %s", conf.SelectedEvents)
	forwarder := nozzle.NewForwarder(
		writerClient, writerEventSerializer,
		conf.SelectedEvents, events, errs, logger,
	)
	err = forwarder.Run(time.Second)
	if err != nil {
		logger.Fatal("Error forwarding", err)
	}
}
