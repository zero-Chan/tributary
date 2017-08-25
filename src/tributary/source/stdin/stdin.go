package stdin

import (
	"bufio"
	"fmt"
	"os"

	"tributary/message"
	"tributary/source"
)

type StdinSource struct {
	source.NormalSource

	reader *bufio.Reader
}

func CreateStdinSource() StdinSource {
	v := StdinSource{
		NormalSource: source.CreateNormalSource(),
	}

	v.SetType("stdin")

	return v
}

func NewStdinSource() *StdinSource {
	v := CreateStdinSource()
	return &v
}

func (s *StdinSource) New() source.Source {
	return NewStdinSource()
}

func (s *StdinSource) StartUp() error {
	if s.Type() != "stdin" {
		return fmt.Errorf("StdinSource StartUp fail: invalid type[%s].", s.Type())
	}

	err := s.NormalSource.StartUp()
	if err != nil {
		return err
	}

	s.reader = bufio.NewReader(os.Stdin)

	go func() {
		// Catch Exception
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("Panic: %v", p)
			}
		}()

		for {
			data, _, err := s.reader.ReadLine()
			if err != nil {
				fmt.Printf("Stdin ReadLine err: %s", err)
			}

			s.Writer.Write(data)
		}
	}()

	return nil
}

func (s *StdinSource) Accept() (message.Message, error) {
	data := <-s.Writer.BufChannel
	return s.Unmarshaler().Unmarshal(data)
}
