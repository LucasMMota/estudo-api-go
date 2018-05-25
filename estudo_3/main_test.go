package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"."
)

var usu main.UsuariosController

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS usuarios (
	id SERIAL,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	CONSTRAINT usuarios_pkey PRIMARY KEY (id)
)`

func TestUsuarios(t *testing.T) {
	usu = main.UsuariosController{}
	usu.Initialize(
		"estudo2",
		"estudo2",
		"estudo2",
	)

	ensureTableExists()

	clearTable()

}

func ensureTableExists() {
	if _, err := usu.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	usu.DB.Exec("DELETE FROM usuarios")
	usu.DB.Exec("ALTER SEQUENCE usuarios_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/usuarios", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	usu.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentUser(t *testing.T) {

}

func TestCreateUsuario(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"Usuario Teste da Silva", "email":"usuario.teste@teste.com.tt"}`)

	req, _ := http.NewRequest("POST", "/usuario", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Usuario Teste da Silva" {
		t.Errorf("Expected user name to be 'Usuario Teste da Silva'. Got '%v'", m["name"])
	}

	if m["price"] != "usuario.teste@teste.com.tt" {
		t.Errorf("Expected user price to be 'usuario.teste@teste.com.tt'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetUsuario(t *testing.T) {
	clearTable()
	addUsuarios(1)

	req, _ := http.NewRequest("GET", "/usuario/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addUsuarios(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		usu.DB.Exec("INSERT INTO usuarios(name, email) VALUES($1, $2)", "Usuario-"+strconv.Itoa(i), "t.tt"+strconv.Itoa(i)+"@t.tt")
	}
}
