package config

import (
	"encoding/json"
	"fmt"
	"log"

	"tributary/source"
	"tributary/source/factory"
)

type Sources map[string]source.Source

func LoadSources(dataSrc map[string]interface{}) (Sources, error) {
	sources := make(Sources)

	for key, val := range dataSrc {
		obj, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Sources[%s] expect an object.", key)
		}

		data, err := json.Marshal(obj)
		if err != nil {
			return nil, fmt.Errorf("Json marshal fail : %v", obj)
		}

		t, exist := obj["type"].(string)
		if !exist {
			return nil, fmt.Errorf("Source[%s] invalid params: not exist type.", key)
		}

		src, err := factory.NewSource(t)
		if err != nil {
			return nil, fmt.Errorf("Create Source fail. reason: %s", err)
		}

		err = json.Unmarshal(data, src)
		if err != nil {
			return nil, fmt.Errorf("json Unmarshal type[%s] fail. data: %s", t, string(data))
		}

		sources[key] = src

		go func(key string) {
			err = sources[key].StartUp()
			if err != nil {
				log.Panicf("Sources[%s] start up fail. reason: %s", key, err)
			}
		}(key)
	}

	return sources, nil
}
