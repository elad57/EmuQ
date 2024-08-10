package broker

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type Broker struct {
	mu    sync.Mutex
	State BrokerState
}

func NewBroker() *Broker {
	return &Broker{
		State: BrokerState{
			Enviorments: make(map[string]Enviorment),
		},
	}
}

func (b *Broker) isQueueExistOnEnviorment(enviorment string, queue string) bool {
	_, isEnviormentExist := b.State.Enviorments[enviorment]
	if !isEnviormentExist {
		return false
	}

	_, isQueueExist := b.State.Enviorments[enviorment].Queues[queue]
	return isQueueExist
}

func (b *Broker) PublishMessage(enviorment string, queue string, message Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.isQueueExistOnEnviorment(enviorment, queue) {
		errorMessage := fmt.Sprintf("'%s:%s' is not exist", enviorment, queue)
		return errors.New(errorMessage)
	}

	for _, subscriber := range(b.State.Enviorments[enviorment].Queues[queue].Subscribers) {
		fmt.Println(subscriber.name)
		subscriber.readFromQueue(message)
	}

	fmt.Printf("Message published to topic '%s:%s': %s\n", enviorment, queue, message.Body)

	return nil
}

func (b *Broker) SubscribeToQueue(enviorment string, queueName string, connection *net.Conn, subscriberName string) ( error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	
	env, ok := b.State.Enviorments[enviorment]
	if !ok {
		errorMessage := fmt.Sprintf("Enviorment '%s' not exist", enviorment)
		return errors.New(errorMessage)
	}
	
	q, ok := env.Queues[queueName]
	
	if !ok {
		errorMessage := fmt.Sprintf("Queue '%s:%s' not exist", enviorment, queueName)
		return errors.New(errorMessage)
	}
	
	subscriber := Consumer{
		connection: connection,
		offset: 0,
		name: subscriberName,
	}
	
	q.Subscribers = append(q.Subscribers, subscriber)
	
	env.Queues[queueName] = q
	
	b.State.Enviorments[enviorment] = env
	fmt.Printf("Subscribed to '%s:%s'\n", enviorment, queueName)
	return nil
}

func (b *Broker) CreateNewEnviorment(enviorment string) error {
	fmt.Printf("Creating '%s' enviorment\n", enviorment)

	_, exists := b.State.Enviorments[enviorment]

	if exists {
		errorMessage := fmt.Sprintf("Enviorment '%s' is already existing", enviorment)
		return errors.New(errorMessage)
	} else {
		b.State.Enviorments[enviorment] = Enviorment{
			Name:   enviorment,
			Queues: make(map[string]Queue),
		}
	}

	fmt.Println(b.State.Enviorments)
	return nil
}

func (b *Broker) CreateNewQueueInEnviorment(queue string, enviorment string) error {
	fmt.Printf("Creating '%s:%s'\n", enviorment, queue)

	_, isEnviormentExist := b.State.Enviorments[enviorment]

	if !isEnviormentExist {
		errorMessage := fmt.Sprintf("Enviorment '%s' not existing", enviorment)
		return errors.New(errorMessage)
	}

	_, isQueueExist := b.State.Enviorments[enviorment].Queues[queue]
	if isQueueExist {
		errorMessage := fmt.Sprintf("Queue '%s' is already existing", queue)
		return errors.New(errorMessage)
	}

	b.State.Enviorments[enviorment].Queues[queue] = Queue{
		Name:        queue,
		Description: "desc",
		CreatedAt:   time.Now(),
	}

	fmt.Println(b.State.Enviorments)
	return nil
}

func (b *Broker) RemoveQueueFromEnviorment(queue string, enviorment string) error {
	fmt.Printf("Removing '%s:%s'\n", enviorment, queue)

	_, isEnviormentExist := b.State.Enviorments[enviorment]
	if !isEnviormentExist {
		errorMessage := fmt.Sprintf("Enviorment '%s' not existing", enviorment)
		return errors.New(errorMessage)
	}

	_, isQueueExist := b.State.Enviorments[enviorment].Queues[queue]
	if !isQueueExist {
		errorMessage := fmt.Sprintf("Queue '%s:%s' not existing", enviorment, queue)
		return errors.New(errorMessage)

	}

	delete(b.State.Enviorments[enviorment].Queues, queue)
	fmt.Println(b.State.Enviorments)
	return nil
}

func (b *Broker) RemoveEnviorment(enviorment string) error {
	fmt.Printf("Removing '%s' enviorment\n", enviorment)
	_, isExist := b.State.Enviorments[enviorment]

	if !isExist {
		errorMessage := fmt.Sprintf("Enviorment '%s' is not existing", enviorment)
		return errors.New(errorMessage)
	}

	delete(b.State.Enviorments, enviorment)
	fmt.Println(b.State.Enviorments)
	return nil
}
