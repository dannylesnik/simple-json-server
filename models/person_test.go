package models

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func createMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestGetPersonShouldReturnData(t *testing.T) {

	db, mock := createMock(t)

	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Firstname", "Lastname"}).AddRow("id1", "Danny", "Lesnik")

	mock.ExpectQuery("^select (.+) from PERSON where*").
		WithArgs("id1").
		WillReturnRows(rows)

	expectedPerson := Person{ID: "id1", Firstname: "Danny", Lastname: "Lesnik"}

	d := &DB{db}

	person, _ := d.GetPerson("id1")

	assert.Equal(t, expectedPerson, *person)

}

func TestGetPersonShouldNotReturnData(t *testing.T) {

	db, mock := createMock(t)

	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Firstname", "Lastname"})

	mock.ExpectQuery("^select (.+) from PERSON where*").
		WithArgs("id2").WillReturnRows(rows)

	d := &DB{db}

	person, err := d.GetPerson("id2")

	assert.Nil(t, person)
	assert.NotNil(t, err)

}

func TestGetPersonShouldThrowException(t *testing.T) {
	db, mock := createMock(t)

	mock.ExpectQuery("^select (.+) from PERSON where*").WithArgs("id2").WillReturnError(errors.New("Some Error"))

	d := &DB{db}

	person, err := d.GetPerson("id2")

	assert.Nil(t, person)
	assert.NotNil(t, err)
}

func TestAddPersonShoulWork(t *testing.T) {
	db, mock := createMock(t)

	defer db.Close()

	mock.ExpectPrepare("^INSERT INTO PERSON*").ExpectExec().WithArgs("id1", "Danny", "Lesnik").WillReturnResult(sqlmock.NewResult(1, 1))

	expectedPerson := Person{ID: "id1", Firstname: "Danny", Lastname: "Lesnik"}

	d := &DB{db}

	actualPerson, err := d.AddPersonToDB(Person{ID: "id1", Firstname: "Danny", Lastname: "Lesnik"})

	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	assert.Equal(t, expectedPerson, *actualPerson)
}

func TestAddPersonThrowsException(t *testing.T) {
	db, mock := createMock(t)

	defer db.Close()

	expectedError := errors.New("Some Error")

	mock.ExpectPrepare("^INSERT INTO PERSON*").ExpectExec().WithArgs("id1", "Danny", "Lesnik").WillReturnError(expectedError)

	d := &DB{db}

	_, err := d.AddPersonToDB(Person{ID: "id1", Firstname: "Danny", Lastname: "Lesnik"})

	assert.Equal(t, expectedError, err)
}

func TestDeletePersonReturnPerson(t *testing.T) {
	db, mock := createMock(t)
	defer db.Close()

	mock.ExpectExec("^delete from PERSON where ID*").WithArgs("id1").WillReturnResult(sqlmock.NewResult(1, 1))

	d := &DB{db}

	result, err := d.DeletePerson("id1")

	if err != nil {
		t.Errorf("error was not expected while deelting person: %s", err)
	}

	var expected int64 = 1
	assert.Equal(t, expected, result)

}

func TestDeletePersonException(t *testing.T) {
	db, mock := createMock(t)
	defer db.Close()

	expectedError := errors.New("Some Error")

	mock.ExpectExec("^delete from PERSON where ID*").WithArgs("id1").WillReturnError(expectedError)

	d := &DB{db}

	_, err := d.DeletePerson("id1")

	assert.Equal(t, expectedError, err)
}

func TestUpdatePersonShouldReturnAffectedRow(t *testing.T) {
	db, mock := createMock(t)
	defer db.Close()

	mock.ExpectExec("^update PERSON set*").WithArgs("Danny", "Lesnik", "id1").WillReturnResult(sqlmock.NewResult(1, 1))

	d := &DB{db}

	result, err := d.UpdatePerson(Person{ID: "id1", Firstname: "Danny", Lastname: "Lesnik"})

	if err != nil {
		t.Errorf("error was not expected while deelting person: %s", err)
	}

	var expected int64 = 1
	assert.Equal(t, expected, result)
}

func TestUpdatePersonShouldThrowException(t *testing.T) {
	db, mock := createMock(t)
	defer db.Close()

	expectedError := errors.New("Some Error")

	mock.ExpectExec("^update PERSON set*").WithArgs("Danny", "Lesnik", "id1").WillReturnError(expectedError)

	d := &DB{db}

	_, err := d.UpdatePerson(Person{ID: "id1", Firstname: "Danny", Lastname: "Lesnik"})

	assert.Equal(t, expectedError, err)

}
