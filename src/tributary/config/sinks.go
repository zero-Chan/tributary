package config

import (
	"encoding/json"
	"fmt"

	"tributary/sink"
	"tributary/sink/factory"
)

type Sinks map[string]sink.Sink

func LoadSinks(dataSinks map[string]interface{}) (Sinks, error) {
	sinks := make(Sinks)

	for key, val := range dataSinks {
		obj, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Sinks[%s] expect an object.", key)
		}

		data, err := json.Marshal(obj)
		if err != nil {
			return nil, fmt.Errorf("Json marshal fail : %v", obj)
		}

		t, exist := obj["type"].(string)
		if !exist {
			return nil, fmt.Errorf("Sink[%s] invalid params: not exist type", key)
		}

		sink, err := factory.NewSink(t)
		if err != nil {
			return nil, fmt.Errorf("Create sink fail. reason: %s", err)
		}

		err = json.Unmarshal(data, sink)
		if err != nil {
			return nil, fmt.Errorf("json Unmarshal type[%s] fail. data: %s", t, string(data))
		}

		sinks[key] = sink

		err = sinks[key].StartUp()
		if err != nil {
			return nil, fmt.Errorf("Sinks[%s] start up fail. reason: %s", key, err)
		}
	}

	return sinks, nil
}
