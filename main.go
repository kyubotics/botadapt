package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/richardchien/botadapt/adapters"
	"github.com/richardchien/botadapt/channels"
	"github.com/richardchien/botadapt/t"
)

func init() {
	// load configuration
	jsonBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Failed to load \"config.json\"")
	}
	json.Unmarshal(jsonBytes, &Config)
	if len(Config.MessageSources) <= 0 {
		log.Fatal("There is no message source defined")
	}
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	messageSourceID, ok := mux.Vars(r)["message_source_id"]
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Incomming event of message source id \"%v\"\n", messageSourceID)

	// construct Event object
	r.ParseForm()
	event := &t.Event{ResponseWriter: w, Request: r, MessageSource: nil}
	for _, messageSource := range Config.MessageSources {
		if id, ok := messageSource["id"]; ok && id == messageSourceID {
			event.MessageSource = messageSource
		}
	}

	if event.MessageSource == nil {
		// there is no such message source id
		http.NotFound(w, r)
		return
	}

	adapter := adapters.Find(event)
	adapter.UnifyEvent(event) // this will start a new goroutine to process the real unifying job
	log.Println("A new goroutine has been started to unify the event")
}

func APIHandler(w http.ResponseWriter, r *http.Request) {

}

func ConsumeEvent() {
	for eventJSON := range channels.EventChan {
		jsonBytes, err := json.Marshal(eventJSON)
		if err == nil {
			go http.Post(Config.PostURL, "application/json", bytes.NewBuffer(jsonBytes))
		}
	}
}

func main() {
	go ConsumeEvent()

	router := mux.NewRouter()
	router.HandleFunc("/event/{message_source_id}", EventHandler).Methods("POST")
	router.HandleFunc("/api", APIHandler).Methods("POST")

	addr := fmt.Sprintf("%v:%v", Config.Host, Config.Port)
	log.Printf("Starting HTTP server at %v...\n", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
