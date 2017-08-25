package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"

	"tributary/message"
	"tributary/source"
)

type TCPSource struct {
	source.NormalSource
	Host      string `json:"host"`
	Port      int    `json:"port"`
	ClientNum int    `json:"client_num"`

	listener net.Listener
}

func CreateTCPSource() TCPSource {
	v := TCPSource{
		NormalSource: source.CreateNormalSource(),
	}

	v.SetType("tcp")

	return v
}

func NewTCPSource() *TCPSource {
	v := CreateTCPSource()
	return &v
}

func (s *TCPSource) New() source.Source {
	return NewTCPSource()
}

func (s *TCPSource) StartUp() error {
	if s.Type() != "tcp" {
		return fmt.Errorf("TCPSource StartUp fail: invalid type[%s].", s.Type())
	}

	err := s.NormalSource.StartUp()
	if err != nil {
		return fmt.Errorf("TCPSource StartUp fail: %s", err)
	}

	s.listener, err = net.Listen("tcp", s.Host+":"+strconv.Itoa(s.Port))
	if err != nil {
		return fmt.Errorf("TCPSource StartUp net.Listen fail: %s", err)
	}
	fmt.Println("Listen: ", s.Host+":"+strconv.Itoa(s.Port))

	go func() {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("Panic: %v", p)
			}
		}()

		for {
			conn, err := s.listener.Accept()
			if err != nil {
				fmt.Printf("TCP Accept err: ", err)
				continue
			}

			// Conn handler
			go func(conn net.Conn) {
				defer conn.Close()
				for {
					// TODO : if client send line : {"a": 1}\n {"b": 2}\n  if will only get the first message.
					line, err := bufio.NewReader(conn).ReadBytes('\n')
					if err != nil {
						if err == io.EOF {
							continue
						}

						fmt.Println("read err : ", err)
						return
					}

					s.Writer.Write(line)
				}
			}(conn)
		}
	}()

	return nil
}

func (s *TCPSource) Accept() (message.Message, error) {
	data := <-s.Writer.BufChannel
	return s.Unmarshaler().Unmarshal([]byte(data))
}
