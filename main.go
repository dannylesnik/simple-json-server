package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dannylesnik/simple-json-server/models"
	"github.com/gorilla/mux"
)

//Error - Error Response Body
type Error struct {
	Error   string `json:"error"`
	Message string `json:"msg"`
	Code    int16  `json:"code"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func isalive(w http.ResponseWriter, r *http.Request) {

	isalive := models.GetIsAliveResponse()
	jsonString, err := json.Marshal(isalive)
	fmt.Println(isalive)
	if err != nil {
		fmt.Println(err)
	} else {
		log.Printf("Your JSON is %s\n", jsonString)
	}
	json.NewEncoder(w).Encode(isalive)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	person, err := models.DeletePerson(eventID)
	if err != nil {
		json.NewEncoder(w).Encode(Error{"Can't get Person", err.Error(), 404})
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(person)
	}
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	person, err := models.GetPerson(eventID)
	if err != nil {
		json.NewEncoder(w).Encode(Error{"Can't get Person", err.Error(), 404})
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(person)
	}
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(Error{"Can't read Request", err.Error(), 400})
	} else {
		if err := models.Unmarshal(reqBody, &person); err != nil {
			json.NewEncoder(w).Encode(Error{"Can't parse JSON Request", err.Error(), 400})
		} else {
			_, err = person.AddPerson()
			if err == nil {
				json.NewEncoder(w).Encode(person)
			} else {
				json.NewEncoder(w).Encode(Error{"Can't can't add new Person", err.Error(), 400})
			}
		}
	}
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(Error{"Can't read Request", err.Error(), 400})
	} else {
		if err := models.Unmarshal(reqBody, &person); err != nil {
			json.NewEncoder(w).Encode(Error{"Can't parse JSON Request", err.Error(), 400})
		} else {
			person.UpdatePerson()
			json.NewEncoder(w).Encode(person)
		}
	}

}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/isalive", isalive)
	router.HandleFunc("/add", createPerson).Methods("POST")
	router.HandleFunc("/persons/{id}", getPerson).Methods("GET")
	router.HandleFunc("/persons/{id}", deletePerson).Methods("DELETE")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(srv)

}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)

}
