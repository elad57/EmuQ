package main

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/elad57/emuq/internal/app"
	"github.com/elad57/emuq/internal/broker"
	// "github.com/elad57/emuq/logger"
)

type Action func(*broker.Broker, app.TCPMessage) error

func createClientMessageToBrokerActionMap() map[app.ClientMessageType]Action {
	return map[app.ClientMessageType]Action{
		app.SUBSCRIBE: func(b *broker.Broker, message app.TCPMessage) error {
			return b.SubscribeToQueue(message.Enviorment, message.Queue, message.Connection, message.Sender)
		},
		app.PRODUCE: func(b *broker.Broker, message app.TCPMessage) error {
			return b.PublishMessage(message.Enviorment, message.Queue, *broker.NewBrokerMessage(message.Payload))
		},
		app.RESPONSE: func(b *broker.Broker, message app.TCPMessage) error { return nil },
	}
}

func main() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Configure zap to write to the file
	fileWriterSyncer := zapcore.AddSync(file)
	consoleSyncer := zapcore.AddSync(os.Stdout)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, fileWriterSyncer, zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, consoleSyncer, zapcore.DebugLevel),
	)
	logger := zap.New(core)

	defer logger.Sync() // flushes buffer, if any

	// logger.Info("This is an info message")
	// logger.Error("This is an error message")
	// logger := logger.NewLogger("./emuq.log", logger.INFO)
	logger.Info("Starting EmuQ...")

	b := broker.NewBroker(logger)

	mapClientMessagesToBrokerActions := createClientMessageToBrokerActionMap()

	tcpServer := app.NewTCPServer(":8082", logger)
	go func() {
		for msg := range tcpServer.MsgChannel {
			msgJSON := string(msg.Payload)
			logger.Sugar().Infof("recevied message %s", msgJSON)
			data := msg
			err := json.Unmarshal([]byte(msgJSON), &data)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			response, _ := json.Marshal(app.TCPMessage{
				Payload:    json.RawMessage([]byte(data.Type)),
				Type:       app.RESPONSE,
				Queue:      msg.Queue,
				Enviorment: msg.Enviorment,
				Sender:     "EmuQ server",
			})
			logger.Debug("writing message")
			(*msg.Connection).Write(response)
			logger.Debug("wrote message")
			action := mapClientMessagesToBrokerActions[data.Type]
			action(b, data)
			logger.Debug("actioned message")
		}
	}()
	go tcpServer.Start()

	httpServer := app.NewHttpServer(b, logger)
	httpServer.StartHttpServer(":8080")
}
