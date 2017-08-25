package processor

import (
	"tributary/message"
)

type Processor interface {
	New() Processor

	SetType(string)
	Type() string

	StartUp() error

	Process(message.Message) (message.Message, error)
}

type NormalProcessor struct {
	ProType string `json:"type"`
}

func CreateNormalProcessor() NormalProcessor {
	v := NormalProcessor{}

	return v
}

func NewNormalProcessor() *NormalProcessor {
	v := CreateNormalProcessor()

	return &v
}

func (p *NormalProcessor) SetType(t string) {
	p.ProType = t
}

func (p NormalProcessor) Type() string {
	return p.ProType
}
