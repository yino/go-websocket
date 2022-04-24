package client

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type Client struct {
	Id           string // client id
	Conn         *websocket.Conn
	readChannel  chan string
	outChannel   chan string
	closeChannel chan byte
	closeOnce    sync.Once
}

func (client *Client) Read() string {
	data, _ := <-client.readChannel
	return data
}
func (client *Client) readLoop() {
	log.Println("readLoop")
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage", err)
			client.Close()
			return
		}
		select {
		case client.readChannel <- string(msg):
			client.Send(string(msg))
		case <-client.closeChannel:
			client.Close()
			return
		}
	}
}

func (client *Client) Send(msg string) {
	client.outChannel <- msg
}
func (client *Client) sendLoop() {
	log.Println("sendLoop")
	for {
		select {
		case data, ok := <-client.outChannel:
			if ok {
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					client.Close()
				}
			}
		case <-client.closeChannel:
			client.Close()
			return
		}
	}
}

func (client *Client) Close() {
	client.Conn.Close()
	client.closeOnce.Do(func() {
		close(client.closeChannel)
	})
	log.Println(client.Id, "close")
}

func (client *Client) heartbeat() {
	for {
		client.Send("heartbeat")
		time.Sleep(1 * time.Second)
	}
}

func NewClient(conn *websocket.Conn) *Client {
	client := &Client{
		Id:           uuid.NewV4().String(),
		Conn:         conn,
		outChannel:   make(chan string, 1000),
		readChannel:  make(chan string, 1000),
		closeChannel: make(chan byte, 1),
	}
	go client.heartbeat()
	go client.readLoop()
	go client.sendLoop()
	return client
}
