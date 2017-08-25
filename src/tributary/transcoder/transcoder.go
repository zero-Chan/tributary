package transcoder

import (
	"tributary/message"
)

type Marshaler interface {
	Marshal(message.Message) (data []byte, err error)
}

type Unmarshaler interface {
	Unmarshal(data []byte) (message.Message, error)
}

type TransCoder interface {
	New() TransCoder
	Marshaler
	Unmarshaler
}
