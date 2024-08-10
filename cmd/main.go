package main

import (
	"encoding/json"
	"fmt"

	"github.com/elad57/emuq/internal/app"
	"github.com/elad57/emuq/internal/broker"
)

type Action func(*broker.Broker, app.TCPMessage) error

func createClientMessageToBrokerActionMap() map[app.ClientMessageType]Action{
	m := map[app.ClientMessageType]Action {
		app.SUBSCRIBE: func(b *broker.Broker, message app.TCPMessage) error { return b.SubscribeToQueue(message.Enviorment, message.Queue, message.Connection, message.Sender) },
		app.PRODUCE: func(b *broker.Broker, message app.TCPMessage) error { return b.PublishMessage(message.Enviorment, message.Queue, *broker.NewBrokerMessage(message.Payload)) },
		app.RESPONSE: func(b *broker.Broker, message app.TCPMessage) error { return nil },
	}

	return m
}

func main() {
	fmt.Println("Starting EmuQ...")

	b := broker.NewBroker()

	mapClientMessagesToBrokerActions := createClientMessageToBrokerActionMap()
		
		
	tcpServer := app.NewTCPServer(":8082")
	go func ()  {
		for msg := range tcpServer.MsgChannel {
			msgJSON := string(msg.Payload)
			fmt.Println("recevied message", msgJSON)
			data := msg
			err:=json.Unmarshal([]byte(msgJSON), &data)
			fmt.Println("decoded message", data)
			if err !=nil {
				fmt.Println(err)
				continue
			}

			response, _ := json.Marshal(app.TCPMessage{
				Payload: json.RawMessage([]byte(data.Type)),
				Type: app.RESPONSE,
				Queue: msg.Queue,
				Enviorment: msg.Enviorment,
				Sender: "EmuQ server",
			})
			(*msg.Connection).Write(response)
			action := mapClientMessagesToBrokerActions[data.Type]
			action(b, data)
		}
		}()
	go tcpServer.Start()
	
	httpServer := app.NewHttpServer(b)
	httpServer.StartHttpServer(":8080")
}