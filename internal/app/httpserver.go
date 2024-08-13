package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elad57/emuq/internal/broker"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)


func NewHttpServer(b *broker.Broker, logger *zap.Logger) *HttpServer {
	router := mux.NewRouter()
	
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Welcome to EmuQ!")
	})	

	router.HandleFunc("/enviorments", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(maps.Keys(b.State.Enviorments))	
	})
	
	router.HandleFunc("/enviorments/{enviorment}", func(w http.ResponseWriter, r *http.Request) {
		enviormentName := mux.Vars(r)["enviorment"]
		b.CreateNewEnviorment(enviormentName)
		json.NewEncoder(w).Encode(maps.Keys(b.State.Enviorments))	
	}).Methods("POST")
	
	router.HandleFunc("/queues/{enviorment}", func(w http.ResponseWriter, r *http.Request) {
		enviormentName := mux.Vars(r)["enviorment"]
		json.NewEncoder(w).Encode(maps.Keys(b.State.Enviorments[enviormentName].Queues))	
	})
	
	router.HandleFunc("/queues/{enviorment}/{queue}", func(w http.ResponseWriter, r *http.Request) {
		enviormentName := mux.Vars(r)["enviorment"]
		queueName := mux.Vars(r)["queue"]
		b.CreateNewQueueInEnviorment(queueName, enviormentName)
		json.NewEncoder(w).Encode(maps.Keys(b.State.Enviorments[enviormentName].Queues))
	}).Methods("POST")
		
	router.HandleFunc("/produce/{enviorment}/{queue}", func(w http.ResponseWriter, r *http.Request) {
		enviormentName := mux.Vars(r)["enviorment"]
		queueName := mux.Vars(r)["queue"]
		body, err := ioutil.ReadAll(r.Body) 
	
		defer r.Body.Close()

		logger.Info(queueName)
		if err != nil {
			logger.Sugar().Error(err)
			w.WriteHeader(500)
			return
		}
		
		b.PublishMessage(enviormentName,queueName, *broker.NewBrokerMessage(body))
		logger.Info(enviormentName)
	}).Methods("POST")
			
	return &HttpServer{
		Broker: b,
		Router: router,
	}
}

func (h *HttpServer)StartHttpServer(port string) {
	log.Fatal(http.ListenAndServe(port, h.Router))
}