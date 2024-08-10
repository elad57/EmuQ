package app

import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/exp/maps"
	"github.com/gorilla/mux"
	"github.com/elad57/emuq/internal/broker"
)


func NewHttpServer(b *broker.Broker) *HttpServer {
	router := mux.NewRouter()
	
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to EmuQ!")
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
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		b.PublishMessage(enviormentName,queueName, *broker.NewBrokerMessage(body))
	}).Methods("POST")
			
	return &HttpServer{
		Broker: b,
		Router: router,
	}
}

func (h *HttpServer)StartHttpServer(port string) {
	log.Fatal(http.ListenAndServe(port, h.Router))
}