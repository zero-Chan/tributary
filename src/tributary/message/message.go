package message

type Message interface {
	Type() string
	ADD(string, interface{})
	GET(string) interface{}
	Copy() Message
}
