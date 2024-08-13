package broker_test

import (
	"net"
	"testing"

	"github.com/elad57/emuq/internal/broker"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type MockConn struct {
	net.Conn
}

type MockLogger struct {
	zap.Logger
}

func TestNewBroker(t *testing.T) {
	broker := broker.NewBroker(zaptest.NewLogger(t))
	if broker == nil {
		t.Error("Expected NewBroker to return a non-nil Broker")
	}
	if len(broker.State.Enviorments) != 0 {
		t.Errorf("Expected no environments, got %d", len(broker.State.Enviorments))
	}
}

func TestCreateNewEnviorment(t *testing.T) {
	broker := broker.NewBroker(zaptest.NewLogger(t))
	err := broker.CreateNewEnviorment("test-env")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, exists := broker.State.Enviorments["test-env"]; !exists {
		t.Error("Expected environment 'test-env' to exist")
	}

	err = broker.CreateNewEnviorment("test-env")
	if err == nil || err.Error() != "Enviorment 'test-env' is already existing" {
		t.Errorf("Expected 'already existing' error, got %v", err)
	}
}

func TestCreateNewQueueInEnviorment(t *testing.T) {
	broker := broker.NewBroker(zaptest.NewLogger(t))
	broker.CreateNewEnviorment("test-env")
	err := broker.CreateNewQueueInEnviorment("test-queue", "test-env")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, exists := broker.State.Enviorments["test-env"].Queues["test-queue"]; !exists {
		t.Error("Expected queue 'test-queue' to exist")
	}

	err = broker.CreateNewQueueInEnviorment("test-queue", "test-env")
	if err == nil || err.Error() != "Queue 'test-queue' is already existing" {
		t.Errorf("Expected 'already existing' error, got %v", err)
	}
}

func TestRemoveQueueFromEnviorment(t *testing.T) {
	broker := broker.NewBroker(zaptest.NewLogger(t))
	broker.CreateNewEnviorment("test-env")
	broker.CreateNewQueueInEnviorment("test-queue", "test-env")
	err := broker.RemoveQueueFromEnviorment("test-queue", "test-env")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, exists := broker.State.Enviorments["test-env"].Queues["test-queue"]; exists {
		t.Error("Expected queue 'test-queue' to be removed")
	}

	err = broker.RemoveQueueFromEnviorment("nonexistent-queue", "test-env")
	if err == nil || err.Error() != "Queue 'test-env:nonexistent-queue' not existing" {
		t.Errorf("Expected 'not existing' error, got %v", err)
	}
}

func TestRemoveEnviorment(t *testing.T) {
	broker := broker.NewBroker(zaptest.NewLogger(t))
	broker.CreateNewEnviorment("test-env")
	err := broker.RemoveEnviorment("test-env")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if _, exists := broker.State.Enviorments["test-env"]; exists {
		t.Error("Expected environment 'test-env' to be removed")
	}

	err = broker.RemoveEnviorment("nonexistent-env")
	if err == nil || err.Error() != "Enviorment 'nonexistent-env' is not existing" {
		t.Errorf("Expected 'not existing' error, got %v", err)
	}
}

func TestPublishMessage(t *testing.T) {
	b := broker.NewBroker(zaptest.NewLogger(t))
	b.CreateNewEnviorment("test-env")
	b.CreateNewQueueInEnviorment("test-queue", "test-env")
	
	message := broker.Message{Body: []byte("hello")}
	err := b.PublishMessage("test-env", "test-queue", message)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSubscribeToQueue(t *testing.T) {
	b := broker.NewBroker(zaptest.NewLogger(t))
	b.CreateNewEnviorment("test-env")
	b.CreateNewQueueInEnviorment("test-queue", "test-env")

	conn := &MockConn{}
	err := b.SubscribeToQueue("test-env", "test-queue", &conn.Conn, "test-subscriber")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(b.State.Enviorments["test-env"].Queues["test-queue"].Subscribers) != 1 {
		t.Errorf("Expected 1 subscriber, got %d", len(b.State.Enviorments["test-env"].Queues["test-queue"].Subscribers))
	}

	err = b.SubscribeToQueue("test-env", "nonexistent-queue", &conn.Conn, "test-subscriber")
	if err == nil || err.Error() != "Queue 'test-env:nonexistent-queue' not exist" {
		t.Errorf("Expected 'not exist' error, got %v", err)
	}
}