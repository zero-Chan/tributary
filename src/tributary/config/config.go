package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"tributary/agent"
)

type Config struct {
	Sources    map[string]interface{} `json:"sources"`
	Processors map[string]interface{} `json:"processors"`
	Sinks      map[string]interface{} `json:"sinks"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "filename", "/etc/barricade/config.json", "Configure filename.")
	flag.StringVar(&configFile, "f", "/etc/barricade/config.json", "Configure filename.")

	flag.Parse()
}

func LoadConfig() (*agent.Agent, error) {
	confFile := configFile

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("Could not find config file.")
	}

	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("Read config file(%s) fail: %v", confFile, err)
	}

	confData := Config{}
	err = json.Unmarshal(data, &confData)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal config fail: %v", err)
	}

	switch {
	case len(confData.Sources) == 0:
		return nil, fmt.Errorf("Must to specify configure[sources]")
	case len(confData.Processors) == 0:
		return nil, fmt.Errorf("Must to specify configure[processors]")
	case len(confData.Sinks) == 0:
		return nil, fmt.Errorf("Must to specify configure[sinks]")
	}

	sources, err := LoadSources(confData.Sources)
	if err != nil {
		return nil, fmt.Errorf("LoadSources fail: %s", err)
	}

	processors, err := LoadProcessors(confData.Processors)
	if err != nil {
		return nil, fmt.Errorf("LoadProcessors fail: %s", err)
	}

	sinks, err := LoadSinks(confData.Sinks)
	if err != nil {
		return nil, fmt.Errorf("LoadSinks fail: %s", err)
	}

	ag, err := agent.NewAgent(sources, processors, sinks)
	if err != nil {
		return nil, fmt.Errorf("New Agent fail: %s", err)
	}

	return ag, nil
}

type Agent struct {
	Sources    Sources
	Processors Processors
	Sinks      Sinks
}

func newService(sources Sources, processors Processors, sinks Sinks) *Agent {
	service := &Agent{
		Sources:    sources,
		Processors: processors,
		Sinks:      sinks,
	}

	return service
}

func (s *Agent) check() error {
	//	if s.Sources == nil || s.Interceptors == nil || s.Processors == nil {
	//		return fmt.Errorf("Service data nil. Sources[%p], Interceptors[%p], processors[%p]\n", s.Sources, s.Interceptors, s.Processors)
	//	}

	//	for _, src := range s.Sources {
	//		for _, itr := range src.Nexts() {
	//			if _, exist := s.Interceptors[itr]; !exist {
	//				return fmt.Errorf("Service Interceptors[%s] not exist.", itr)
	//			}
	//		}
	//	}

	//	for _, itr := range s.Interceptors {
	//		for _, prs := range itr.Nexts() {
	//			if _, exist := s.Processors[prs]; !exist {
	//				return fmt.Errorf("Service processors[%s] not exist.", prs)
	//			}
	//		}
	//	}

	//	for _, pro := range s.Processors {
	//		for _, snk := range pro.Nexts() {
	//			if _, exist := s.Sinks[snk]; !exist {
	//				return fmt.Errorf("Service sinks[%s] not exist.", snk)
	//			}
	//		}
	//	}

	return nil
}

func (s *Agent) Start() error {
	err := s.check()
	if err != nil {
		return err
	}
	// for () { go Sources.Accept() }

	return nil
}
