#!/usr/bin/env bash

export NOZZLE_UAA_URL=https://uaa.bosh-lite.com
# this uaa client credential set isn't built in; see README.md
export NOZZLE_USERNAME=my-nozzle-client
export NOZZLE_PASSWORD=password123
export NOZZLE_TRAFFIC_CONTROLLER_URL=wss://doppler.bosh-lite.com:4443
export NOZZLE_FIREHOSE_SUBSCRIPTION_ID=my-firehose-subscription
export NOZZLE_SKIP_SSL=true

export EVENTS_ALL=HttpStartStop,LogMessage,ValueMetric,CounterEvent,Error,ContainerMetric
export EVENTS_PLATFORM=HttpStartStop,ValueMetric,CounterEvent,Error
export EVENTS_APP=HttpStartStop,LogMessage,ContainerMetric
export NOZZLE_SELECTED_EVENTS=${EVENTS_APP}

go run main.go
