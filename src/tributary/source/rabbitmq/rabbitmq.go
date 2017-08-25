package rabbitmq

import (
	"fmt"
	//	"time"

	"tributary/message"
	"tributary/message/json"
	"tributary/source"
)

type RabbitmqSource struct {
	source.NormalSource
	// rabbitmqConsumer
}

func CreateRabbitmqSource() RabbitmqSource {
	v := RabbitmqSource{
		NormalSource: source.CreateNormalSource(),
	}

	v.SetType("rabbitmq")

	return v
}

func NewRabbitmqSource() *RabbitmqSource {
	v := CreateRabbitmqSource()
	return &v
}

func (mq *RabbitmqSource) New() source.Source {
	return NewRabbitmqSource()
}

func (mq *RabbitmqSource) StartUp() error {
	if mq.Type() != "rabbitmq" {
		return fmt.Errorf("RabbitmqSource StartUp fail: invalid type[%s].", mq.Type())
	}

	// TODO
	fmt.Println("RabbitmqSource start up.")

	return nil
}

func (mq *RabbitmqSource) Accept() (message.Message, error) {

	// TODO
	m := json.NewJsonMessage()
	m.ADD("key1", 1)
	m.ADD("key2", "hello")

	fmt.Println("Rabbitmq accept message: ", m)

	//	time.Sleep(1e9)

	return m, nil
}
