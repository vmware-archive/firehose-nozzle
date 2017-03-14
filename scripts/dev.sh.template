#!/usr/bin/env bash

auth_method="api"
#auth_method="uaa"

if [[ "${auth_method}" == "api" ]]; then
    echo "Doing API auth"
    export NOZZLE_API_URL=https://api.bosh-lite.com
    export NOZZLE_USERNAME=admin
    export NOZZLE_PASSWORD=admin #admin / admin is the bosh-lite credential set
else
    echo "Doing UAA auth"
    export NOZZLE_UAA_URL=https://uaa.bosh-lite.com
    export NOZZLE_TRAFFIC_CONTROLLER_URL=wss://doppler.bosh-lite.com:4443
    export NOZZLE_USERNAME=<username>
    export NOZZLE_PASSWORD=<password>
fi

export NOZZLE_FIREHOSE_SUBSCRIPTION_ID=firehose-subscription-id
export NOZZLE_SKIP_SSL=true
export NOZZLE_SELECTED_EVENTS=ValueMetric,CounterEvent

go run main.go
