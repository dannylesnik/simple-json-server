package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//DB - structure that is holding the Database
type DB struct {
	*sql.DB
}

//Datastore - Generic Database unterface
type Datastore interface {
	AddPersonToDB(person Person) (*Person, error)
	GetPerson(id string) (*Person, error)
	DeletePerson(id string) (int64, error)
	UpdatePerson(person Person) (int64, error)
}

//InitDB - creates new =connection instance
func InitDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
