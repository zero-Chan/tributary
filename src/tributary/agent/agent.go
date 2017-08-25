package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tributary/message"
	"tributary/processor"
	"tributary/sink"
	"tributary/source"
)

func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		log.Printf("comet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			return
		case syscall.SIGHUP:
			// TODO : reload()
		default:
			return
		}
	}
}

type Agent struct {
	sources    map[string]source.Source
	processors map[string]processor.Processor
	sinks      map[string]sink.Sink

	proChannelMgr  map[string]processChannelManager
	sinkChannelMgr map[string]sinkChannelManager
}

func CreateAgent(sources map[string]source.Source, processors map[string]processor.Processor, sinks map[string]sink.Sink) (Agent, error) {
	v := Agent{}

	switch {
	case sources == nil:
		return v, fmt.Errorf("Invalid: sources is nil.")
	case processors == nil:
		return v, fmt.Errorf("Invalid: processors is nil.")
	case sinks == nil:
		return v, fmt.Errorf("Invalid: sinks is nil.")
	}

	v.sources = sources
	v.processors = processors
	v.sinks = sinks
	v.proChannelMgr = make(map[string]processChannelManager)
	v.sinkChannelMgr = make(map[string]sinkChannelManager)

	return v, nil
}

func NewAgent(sources map[string]source.Source, processors map[string]processor.Processor, sinks map[string]sink.Sink) (*Agent, error) {
	v, err := CreateAgent(sources, processors, sinks)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (ag *Agent) Start() {
	ag.linkProcessors()
	ag.createProcessChannel()
	ag.createSources()
	ag.createSinks()
	ag.distributer()
	//	ag.monitor()

	select {}
}

func (ag *Agent) linkProcessors() {
	// Link Sources->Processors
	for srcKey, src := range ag.sources {
		for _, proKey := range src.ProcessorKeys() {
			pro, exist := ag.processors[proKey]
			if !exist {
				log.Panicf("Sources[%s] import not exist processor[%s]", srcKey, proKey)
			}

			src.SetProcessor(pro)
		}
	}

	// Link Sinks->Processors
	for sinkKey, sk := range ag.sinks {
		for _, proKey := range sk.ProcessorKeys() {
			pro, exist := ag.processors[proKey]
			if !exist {
				log.Panicf("Sinks[%s] import not exist processor[%s]", sinkKey, proKey)
			}

			sk.SetProcessor(pro)
		}
	}
}

// createProcessChannel use to create process handler, and make mgr to record in/out channel
func (ag *Agent) createProcessChannel() {
	for key, pro := range ag.processors {
		in := make(chan message.Message, len(ag.sources))
		out := make(chan message.Message, len(ag.sinks))
		ag.proChannelMgr[key] = processChannelManager{
			in:   in,
			out:  out,
			name: key,
			pro:  pro,
		}

		go func(in <-chan message.Message, out chan<- message.Message, pro processor.Processor, name string) {
			for inMsg := range in {

				outMsg, err := func() (out message.Message, err error) {
					// Catch Exception
					defer func() {
						if p := recover(); p != nil {
							err = fmt.Errorf("Panic: %v", p)
						}
					}()

					return pro.Process(inMsg)
				}()

				if err != nil {
					log.Printf("Processor[%s] Handle error: %s", name, err)
					continue
				}

				out <- outMsg
			}

		}(in, out, pro, key)
	}
}

// createSources use to create source handler, and publish msg entity to processMgr in channel
func (ag *Agent) createSources() {
	for key, src := range ag.sources {

		go func(src source.Source, name string) {
			for {
				msg, err := func() (msg message.Message, err error) {
					// Catch Exception
					defer func() {
						if p := recover(); p != nil {
							err = fmt.Errorf("Panic: %v", p)
						}
					}()

					return src.Accept()
				}()

				if err != nil {
					log.Printf("Source[%s] Accept error: %s", name, err)
					continue
				}

				for _, proKey := range src.ProcessorKeys() {
					ag.proChannelMgr[proKey].in <- msg.Copy()
				}
			}
		}(src, key)
	}
}

// createSources use to create sink handler, and make mgr to record in channel
func (ag *Agent) createSinks() {
	for key, snk := range ag.sinks {
		in := make(chan message.Message, len(ag.processors))
		ag.sinkChannelMgr[key] = sinkChannelManager{
			in:   in,
			name: key,
			snk:  snk,
		}

		go func(in <-chan message.Message, snk sink.Sink, name string) {
			for msg := range in {
				err := func() (err error) {
					// Catch Exception
					defer func() {
						if p := recover(); p != nil {
							err = fmt.Errorf("Panic: %v", p)
						}
					}()

					return snk.Publish(msg.Copy())
				}()

				if err != nil {
					log.Printf("Sink[%s] Publish error: %s", name, err)
					continue
				}
			}

		}(in, snk, key)

	}
}

// distributer use to get msg from processorMgr out channel, and distributer msg entity to every sinkMgr in channel
func (ag *Agent) distributer() {
	for proKey, proMgr := range ag.proChannelMgr {
		sinkLIst := make([]string, 0)

		for sinkKey, snk := range ag.sinks {
			for _, snkProKey := range snk.ProcessorKeys() {
				if snkProKey == proKey {
					sinkLIst = append(sinkLIst, sinkKey)
				}
			}
		}

		go func(proKey string, proMgr processChannelManager, sinkLIst []string, sinkMgr map[string]sinkChannelManager) {
			for msg := range proMgr.out {

				for _, sinkKey := range sinkLIst {
					sinkMgr[sinkKey].in <- msg.Copy()
				}
			}
		}(proKey, proMgr, sinkLIst, ag.sinkChannelMgr)
	}

}

func (ag *Agent) monitor() {
	go func(ag *Agent) {
		for {
			var logProBuf string
			for proKey, proMgr := range ag.proChannelMgr {
				logProBuf += fmt.Sprintf("ProcessChannel: [%s] in: (%d/%d). out(%d/%d)      ", proKey, len(proMgr.in), cap(proMgr.in), len(proMgr.out), cap(proMgr.out))
			}

			var logSinkBuf string
			for sinkKey, sinkMgr := range ag.sinkChannelMgr {
				logSinkBuf += fmt.Sprintf("SinkChannel: [%s] in: (%d/%d)		", sinkKey, len(sinkMgr.in), cap(sinkMgr.in))
			}

			// Print

			fmt.Printf("=================== Monitor ===================\n")
			log.Printf("%s\n", logProBuf)
			log.Printf("%s\n", logSinkBuf)
			fmt.Printf("=============================================\n")

			time.Sleep(10e9)
		}

	}(ag)
}

type processChannelManager struct {
	in   chan message.Message
	out  chan message.Message
	name string
	pro  processor.Processor
}

type sinkChannelManager struct {
	in   chan message.Message
	name string
	snk  sink.Sink
}
