package factory

import (
	"fmt"

	"tributary/source"
	"tributary/source/rabbitmq"
	"tributary/source/stdin"
	"tributary/source/tcp"
)

type sourceFactory map[string]source.Source

var globalSourceFactory = make(sourceFactory)

func init() {
	globalSourceFactory.add("rabbitmq", rabbitmq.NewRabbitmqSource())
	globalSourceFactory.add("stdin", stdin.NewStdinSource())
	globalSourceFactory.add("tcp", tcp.NewTCPSource())
}

func (sf *sourceFactory) add(t string, s source.Source) {
	(*sf)[t] = s
}

func (sf *sourceFactory) remove(t string) {
	delete(*sf, t)
}

func (sf *sourceFactory) get(t string) (source.Source, error) {
	s, exists := (*sf)[t]
	if !exists {
		return nil, fmt.Errorf("sourceFactory does not exist type[%s]", t)
	}

	return s.New(), nil
}

func (sf sourceFactory) list() []string {
	res := make([]string, len(sf))
	for k, _ := range sf {
		res = append(res, k)
	}

	return res
}

func AddSourceFactory(t string, s source.Source) {
	globalSourceFactory.add(t, s)
}

func NewSource(t string) (source.Source, error) {
	return globalSourceFactory.get(t)
}

func ListSourceFactory() []string {
	return globalSourceFactory.list()
}

func DeleteSourceFactory(t string) {
	globalSourceFactory.remove(t)
}
