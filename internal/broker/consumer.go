package broker

import (
	"encoding/json"
	"fmt"
)

func (consumer Consumer) readFromQueue(message Message) {
	conn := *consumer.connection
	jsonMessage, err := json.Marshal(message.Body)
	fmt.Println(consumer.name, "has read from queue message:", jsonMessage)
	if err != nil {	
		fmt.Println("error", err)
	} else {
		conn.Write(message.Body)
	}
}