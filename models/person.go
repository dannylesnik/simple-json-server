package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

//Person - This is person type
type Person struct {
	ID        string `json:"id" binding:"exists"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
}

var persons = make(map[string]Person)

//AddPerson - adds new Person
func (person Person) AddPerson() (Person, error) {

	if _, ok := persons[person.ID]; ok {
		return person, errors.New("Person with this ID already exists")
	}
	persons[person.ID] = person
	fmt.Println("person added")
	return person, nil
}

//GetPerson - get single person
func GetPerson(id string) (*Person, error) {
	if person, ok := persons[id]; ok {
		return &person, nil
	} else {
		return nil, errors.New("person not Found")
	}
}

//DeletePerson - delete single person
func DeletePerson(id string) (*Person, error) {
	if person, ok := persons[id]; ok {
		delete(persons, id)
		fmt.Println("deleting object")
		return &person, nil
	} else {
		return nil, errors.New("person not Found")
	}
}

//UpdatePerson - 
func (person Person) UpdatePerson() Person {


		persons[person.ID] = person
	return person
}
//Unmarshal --- unmarshal person
func Unmarshal(data []byte, person *Person) error {
	err := json.Unmarshal(data, &person)

	if err != nil {
		return err
	}
	if person.ID == "" {
		return errors.New("json: ID is missing")
	}
	if person.Firstname == "" {
		return errors.New("json: Firstname is missing")
	}
	if person.Lastname == "" {
		return errors.New("json: Lastname is missing")
	}
	return nil
}
