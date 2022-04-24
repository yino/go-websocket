package client

type Conn interface {
	Read() string
	Send(msg string)
	Close()
}
