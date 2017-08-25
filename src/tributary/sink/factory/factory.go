package factory

import (
	"fmt"

	"tributary/sink"
	"tributary/sink/rabbitmq"
	"tributary/sink/stdout"
)

type sinkFactory map[string]sink.Sink

var globalSinkFactory = make(sinkFactory)

func init() {
	globalSinkFactory.add("rabbitmq", rabbitmq.NewRabbitmqSink())
	globalSinkFactory.add("stdout", stdout.NewStdoutSink())
}

func (sf *sinkFactory) add(t string, s sink.Sink) {
	(*sf)[t] = s
}

func (sf *sinkFactory) remove(t string) {
	delete(*sf, t)
}

func (sf *sinkFactory) get(t string) (sink.Sink, error) {
	p, exists := (*sf)[t]
	if !exists {
		return nil, fmt.Errorf("sinkFactory does not exist type[%s]", t)
	}

	return p.New(), nil
}

func (sf sinkFactory) list() []string {
	res := make([]string, len(sf))
	for k, _ := range sf {
		res = append(res, k)
	}

	return res
}

func AddSinkFactory(t string, s sink.Sink) {
	globalSinkFactory.add(t, s)
}

func NewSink(t string) (sink.Sink, error) {
	return globalSinkFactory.get(t)
}

func ListSinkFactory() []string {
	return globalSinkFactory.list()
}

func DeleteSinkFactory(t string) {
	globalSinkFactory.remove(t)
}
