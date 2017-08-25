package config

import (
	"encoding/json"
	"fmt"

	"tributary/processor"
	"tributary/processor/factory"
)

type Processors map[string]processor.Processor

func LoadProcessors(dataPro map[string]interface{}) (Processors, error) {
	processors := make(Processors)

	for key, val := range dataPro {
		obj, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Processors[%s] expect an object.", key)
		}

		data, err := json.Marshal(obj)
		if err != nil {
			return nil, fmt.Errorf("Json marshal fail : %v", obj)
		}

		t, exist := obj["type"].(string)
		if !exist {
			return nil, fmt.Errorf("Porcessor[%s] invalid params: not exist type.", key)
		}

		pro, err := factory.NewProcessor(t)
		if err != nil {
			return nil, fmt.Errorf("Create processor fail. reason: %s", err)
		}

		err = json.Unmarshal(data, pro)
		if err != nil {
			return nil, fmt.Errorf("json Unmarshal type[%s] fail. data: %s", t, string(data))
		}

		processors[key] = pro

		err = processors[key].StartUp()
		if err != nil {
			return nil, fmt.Errorf("Processors[%s] start up fail. reason: %s", key, err)
		}
	}

	return processors, nil
}
