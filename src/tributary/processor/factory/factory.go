package factory

import (
	"fmt"

	"tributary/processor"
	"tributary/processor/logger_alerter"
	"tributary/processor/timestamp"
)

type processorFactory map[string]processor.Processor

var globalProcessorFactory = make(processorFactory)

func init() {
	globalProcessorFactory.add("loggerAlerter", logger_alerter.NewLoggerAlerterProcessor())
	globalProcessorFactory.add("timestamp", timestamp.NewTimestampProcessor())
}

func (pf *processorFactory) add(t string, p processor.Processor) {
	(*pf)[t] = p
}

func (pf *processorFactory) remove(t string) {
	delete(*pf, t)
}

func (pf *processorFactory) get(t string) (processor.Processor, error) {
	p, exists := (*pf)[t]
	if !exists {
		return nil, fmt.Errorf("processorFactory does not exist type[%s]", t)
	}

	return p.New(), nil
}

func (pf processorFactory) list() []string {
	res := make([]string, len(pf))
	for k, _ := range pf {
		res = append(res, k)
	}

	return res
}

func AddProcessorFactory(t string, p processor.Processor) {
	globalProcessorFactory.add(t, p)
}

func NewProcessor(t string) (processor.Processor, error) {
	return globalProcessorFactory.get(t)
}

func ListProcessorFactory() []string {
	return globalProcessorFactory.list()
}

func DeleteProcessorFactory(t string) {
	globalProcessorFactory.remove(t)
}
