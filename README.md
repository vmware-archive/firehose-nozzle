## Firehose Nozzle

A minimal example for connecting to Cloud Foundry's
[Loggregator](https://github.com/cloudfoundry/loggregator)
system.

### Setup
This code uses [Glide](https://glide.sh/) for dependency management.
After you [install Glide](https://glide.sh/), run `glide install` to
install the Firehose Nozzle dependencies. After that, a `go install`
should build and install the `firehose-nozzle` executable.

Ensure that you have installed firehose-nozzle using `go get github.com/cf-platform-eng/firehose-nozzle`,
this will put the dependent packages in the right place.

There are two options for creating credentials that can talk to the API:
* UAA API user account (choose this route unless you have a reason not to)
* UAA Client

#### Option 1: UAA API User Account

Create a UAA user with access to the Firehose and Cloud Controller:

* Install the UAA CLI, uaac.
```
gem install cf-uaac
```

* Use the uaac target uaa.YOUR-SYSTEM-DOMAIN command to target your UAA server.
```
uaac target uaa.sys.example.com
```

* Record the uaa:admin:client_secret from either
    * The cf deployment manifest or
    *  The Ops Manager UI, click on Pivotal Elastic Runtime tile, go to credentials tab, UAA -> Admin Client Credentials

* Authenticate to UAA using the client_secret from the previous step
```
uaac token client get admin -s ADMIN-CLIENT-SECRET
```

* Create a Nozzle user for your app with the password of your choosing.
```
uaac -t user add my-firehose-nozzle-user --password PASSWORD --emails na
```

* Add the user to the Cloud Controller Admin Read-Only group.
```
uaac -t member add cloud_controller.admin_read_only my-firehose-nozzle-user
```

* Add the user to the Doppler Firehose group.
```
uaac -t member add doppler.firehose my-firehose-nozzle-user
```


#### Option 2: UAA Client

Create a UAA client with access to the Firehose and Cloud Controller:

* Install the UAA CLI, uaac.
```
gem install cf-uaac
```

* Use the uaac target uaa.YOUR-SYSTEM-DOMAIN command to target your UAA server.
```
uaac target uaa.sys.example.com
```

* Record the uaa:admin:client_secret from either
    * The cf deployment manifest or
    *  The Ops Manager UI, click on Pivotal Elastic Runtime tile, go to credentials tab, UAA -> Admin Client Credentials

* Authenticate to UAA using the client_secret from the previous step
```
uaac token client get admin -s ADMIN-CLIENT-SECRET
```

* Create a new UAA client for your firehose
```
uaac client add my-firehose-nozzle \
    --access_token_validity 1209600 \
    --authorized_grant_types authorization_code,client_credentials,refresh_token \
    -s <SECRET> \
    --scope openid,oauth.approvals,doppler.firehose \
    --authorities oauth.login,doppler.firehose
```

For information about creating a UAA user, see the [Creating and Managing Users with the UAA CLI](http://docs.pivotal.io/pivotalcf/adminguide/uaa-user-management.html) topic.

### Development

For development against
[bosh-lite](https://github.com/cloudfoundry/bosh-lite),
copy `scripts/dev.sh.template` to `scripts/dev.sh` and supply missing values.
Then run `./scripts/dev.sh` to see events on standard out.

Install dependencies
```bash
glide install
```

Setup Tests
```bash
go get github.com/onsi/ginkgo/ginkgo  # installs the ginkgo CLI
go get github.com/onsi/gomega         # fetches the matcher library
```

Run test
from toplevel directory
```bash
run ginkgo -r  -skipPackage vendor/   # runs test recursively
```


### References

Other nozzles
* https://github.com/cloudfoundry-incubator/datadog-firehose-nozzle
* https://github.com/cloudfoundry-incubator/datadog-firehose-nozzle-release
* https://github.com/cloudfoundry-community/splunk-firehose-nozzle
* https://github.com/cloudfoundry-community/splunk-firehose-nozzle-release
* https://github.com/cloudfoundry-community/firehose-to-syslog
* https://github.com/cloudfoundry/firehose-plugin
* https://github.com/rakutentech/kafka-firehose-nozzle
* https://github.com/pivotal-cf/graphite-nozzle

General
* https://github.com/cloudfoundry/dropsonde-protocol/tree/master/events
