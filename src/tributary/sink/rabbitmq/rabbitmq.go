package rabbitmq

import (
	"fmt"
	//	"time"

	"tributary/message"
	"tributary/sink"
)

type RabbitmqSink struct {
	sink.NormalSink
	// rabbitmqProducer
}

func CreateRabbitmqSink() RabbitmqSink {
	v := RabbitmqSink{
		NormalSink: sink.CreateNormalSink(),
	}

	v.SetType("rabbitmq")

	return v
}

func NewRabbitmqSink() *RabbitmqSink {
	v := CreateRabbitmqSink()
	return &v
}

func (mq *RabbitmqSink) New() sink.Sink {
	return NewRabbitmqSink()
}

func (mq *RabbitmqSink) StartUp() error {
	if mq.Type() != "rabbitmq" {
		return fmt.Errorf("RabbitmqSource StartUp fail: invalid type[%s].", mq.Type())
	}

	// TODO
	fmt.Println("RabbitmqSink start up.")

	return nil
}

func (mq *RabbitmqSink) Publish(msg message.Message) error {

	// TODO
	fmt.Println("rabbtimq publish: ", msg)

	//	time.Sleep(1e9)

	return nil
}
