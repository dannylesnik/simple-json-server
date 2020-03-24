package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dannylesnik/simple-json-server/models"
	"github.com/gorilla/mux"
)

//Env - Dependencies Injection
type Env struct {
	DB models.Datastore
}

//GetPerson -
// GetPerson godoc
// @Summary Get Person by ID
// @Description Get Person record as JSON by Person's ID
// @Tags Person API
// @ID GetPerson
// @Produce  json
// @Param id path string true "Person ID"
// @Success 200 {object} models.Person
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/v1/persons/{id} [get]
func (env *Env) GetPerson(w http.ResponseWriter, r *http.Request) {
	personID := mux.Vars(r)["id"]
	log.Printf(" Person ID %s", personID)

	person, err := env.DB.GetPerson(personID)
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(models.Error{"Can't get Person", err.Error(), 404})
	} else if err != nil {
		json.NewEncoder(w).Encode(models.Error{"Can't get Person", err.Error(), 500})
	} else {
		json.NewEncoder(w).Encode(person)
	}
}

//UpdatePerson --
// UpdatePerson godoc
// @Summary Update Person record
// @Description Update Person record
// @Tags Person API
// @ID UpdatePerson
// @Accept  json
// @Produce  json
// @Param person body models.Person true "Person Record as JSON"
// @Success 200 {object} models.Person
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/v1/update [put]
func (env *Env) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(models.Error{"Can't read Request", err.Error(), 400})
	} else {
		if err := models.Unmarshal(reqBody, &person); err != nil {
			json.NewEncoder(w).Encode(models.Error{"Can't parse JSON Request", err.Error(), 400})
		} else {
			result, err := env.DB.UpdatePerson(person)
			if err != nil {
				json.NewEncoder(w).Encode(models.Error{"Can't Update Person!!", err.Error(), 500})
			} else if result == 0 {
				json.NewEncoder(w).Encode(models.Error{"Person with such ID doesnt exist!!!", errors.New("Query returned 0 affected records").Error(), 404})
			} else {
				json.NewEncoder(w).Encode(person)
			}
		}
	}
}

//CreatePerson -
// CreatePerson godoc
// @Summary Create new Person record
// @Description Create new Person record
// @Tags Person API
// @ID CreatePerson
// @Accept  json
// @Produce  json
// @Param person body models.Person true "Person Record as JSON"
// @Success 201 {object} models.Person
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/v1/add [post]
func (env *Env) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(models.Error{"Can't read Request", err.Error(), 400})
	} else {
		if err := models.Unmarshal(reqBody, &person); err != nil {
			json.NewEncoder(w).Encode(models.Error{"Can't parse JSON Request", err.Error(), 400})
		} else {
			_, err := env.DB.AddPersonToDB(person)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(models.Error{"Can't Save to DB!!", err.Error(), 500})
			} else {
				json.NewEncoder(w).Encode(person)
			}
		}
	}
}

//DeletePerson -
// DeletePerson godoc
// @Summary Delete Person by ID
// @Description Delete Person record by it's ID
// @Tags Person API
// @ID DeletePerson
// @Produce  json
// @Param id path string true "Person ID"
// @Success 200 {object} models.Person
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /api/v1/persons/{id} [delete]
func (env *Env) DeletePerson(w http.ResponseWriter, r *http.Request) {
	personID := mux.Vars(r)["id"]
	log.Printf(" Event ID %s", personID)
	person, err := env.DB.GetPerson(personID)
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(models.Error{"Can't get Person", err.Error(), 404})
	} else if err != nil {
		json.NewEncoder(w).Encode(models.Error{"Can't delete Person", err.Error(), 500})
	} else {
		result, err := env.DB.DeletePerson(personID)
		if err != nil {
			json.NewEncoder(w).Encode(models.Error{"Can't delete Person", err.Error(), 500})
		} else if result == 0 {
			json.NewEncoder(w).Encode(models.Error{"Can't delete Person", "Person does not exist!!!", 404})
		} else {
			json.NewEncoder(w).Encode(person)
		}
	}
}
