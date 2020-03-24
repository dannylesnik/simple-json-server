package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/dannylesnik/simple-json-server/docs"
	"github.com/dannylesnik/simple-json-server/handlers"
	"github.com/dannylesnik/simple-json-server/models"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
}

// IsAlive godoc
// @Summary Liveness
// @Description Returns hostname, IP, time and Version of current service
// @Tags Liveness API
// @Produce  json
// @Success 200 {object} models.IsAliveResponse
// @Router /api/v1/isalive [get]
func isalive(w http.ResponseWriter, r *http.Request) {
	isalive := models.GetIsAliveResponse()
	json.NewEncoder(w).Encode(isalive)
}

// @title Simple Json Server
// @version 1.0
// @description This is REST API for Simple-json-server demo application.
// @termsOfService http://hithub.com/dannylesnik/

// @contact.name Danny Lesnik
// @contact.url http://hithub.com/dannylesnik/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @x-extension-openapi {"example": "value on a json format"}
func main() {

	dbase, err := models.InitDB("root:my-said2000@/test")
	if err != nil {
		log.Panic(err)
	}

	env := &handlers.Env{DB: dbase}

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	router.HandleFunc("/", homeLink)
	s := router.PathPrefix("/api/v1").Subrouter()

	s.HandleFunc("/isalive", isalive)
	s.HandleFunc("/add", env.CreatePerson).Methods("POST")
	s.HandleFunc("/persons/{id}", env.GetPerson).Methods("GET")
	s.HandleFunc("/persons/{id}", env.DeletePerson).Methods("DELETE")
	s.HandleFunc("/update", env.UpdatePerson).Methods("PUT")

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
