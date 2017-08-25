package factory

import (
	"fmt"

	"tributary/transcoder"
	"tributary/transcoder/json"
)

type transcoderFactory map[string]transcoder.TransCoder

var globalTranscoderFactory = make(transcoderFactory)

func init() {
	globalTranscoderFactory.add("json", json.NewJsonTransCoder())
}

func (tf *transcoderFactory) add(t string, s transcoder.TransCoder) {
	(*tf)[t] = s
}

func (tf *transcoderFactory) remove(t string) {
	delete(*tf, t)
}

func (tf *transcoderFactory) get(t string) (transcoder.TransCoder, error) {
	s, exists := (*tf)[t]
	if !exists {
		return nil, fmt.Errorf("transcoderFactory does not exist type[%s]", t)
	}

	return s.New(), nil
}

func (tf transcoderFactory) list() []string {
	res := make([]string, 0)
	for k, _ := range tf {
		res = append(res, k)
	}

	return res
}

func AddTranscoderFactory(t string, s transcoder.TransCoder) {
	globalTranscoderFactory.add(t, s)
}

func CreateTranscoder(t string) (transcoder.TransCoder, error) {
	return globalTranscoderFactory.get(t)
}

func ListTranscoderFactory() []string {
	return globalTranscoderFactory.list()
}

func DeleteTranscoderFactory(t string) {
	globalTranscoderFactory.remove(t)
}
