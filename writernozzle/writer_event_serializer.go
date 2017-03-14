package writernozzle

import (
	"fmt"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/cloudfoundry-community/go-cfclient"
)

type WriterEventSerializer struct{}

func NewWriterEventSerializer() *WriterEventSerializer {
	return &WriterEventSerializer{}
}

func (w *WriterEventSerializer) BuildHttpStartStopEvent(event *events.Envelope, app *cfclient.App) interface{} {
	return genericSerializer(event, app)
}

func (w *WriterEventSerializer) BuildLogMessageEvent(event *events.Envelope, app *cfclient.App) interface{} {
	return genericSerializer(event, app)
}

func (w *WriterEventSerializer) BuildValueMetricEvent(event *events.Envelope) interface{} {
	return genericSerializer(event, nil)
}

func (w *WriterEventSerializer) BuildCounterEvent(event *events.Envelope) interface{} {
	return genericSerializer(event, nil)
}

func (w *WriterEventSerializer) BuildErrorEvent(event *events.Envelope) interface{} {
	return genericSerializer(event, nil)
}

func (w *WriterEventSerializer) BuildContainerEvent(event *events.Envelope, app *cfclient.App) interface{} {
	return genericSerializer(event, app)
}

func genericSerializer(event *events.Envelope, app *cfclient.App) []byte {
	return []byte(fmt.Sprintf("%+v %+v\n\n", event, app))
}
