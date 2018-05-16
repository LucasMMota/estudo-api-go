package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//App will need two methods that initialize and run the application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//Initialize method will take in the details required to connect to the database
func (a *App) Initialize(user, password, dbname string) {
	var err error
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	// fmt.Println(connectionString)
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
}

//Run method will simply start the application.
func (a *App) Run(addr string) {}
