package usuarios

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type UsuariosController struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize Inicia o módulo de usuários fazendo a conexão com o banco e criando as rotas
func (uc *UsuariosController) Initialize(user, password, dbname string) {
	var err error

	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	uc.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	uc.Router = mux.NewRouter()
	uc.initializeRoutes()
}

//Cria as rotas
func (uc *UsuariosController) initializeRoutes() {
	uc.Router.HandleFunc("/usuarios", uc.getUsuarios).Methods("GET")
	uc.Router.HandleFunc("/usuario", uc.createUsuario).Methods("POST")
	uc.Router.HandleFunc("/usuario/{id:[0-9]+}", uc.getUsuario).Methods("GET")
	uc.Router.HandleFunc("/usuario/{id:[0-9]+}", uc.updateUsuario).Methods("PUT")
	uc.Router.HandleFunc("/usuario/{id:[0-9]+}", uc.deleteUsuario).Methods("DELETE")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (uc *UsuariosController) getUsuarios(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	usuarios, err := getUsuarios(uc.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, usuarios)
}

func (uc *UsuariosController) createUsuario(w http.ResponseWriter, r *http.Request) {
	var usu usuario
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usu); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close() //

	if err := usu.createUsuario(uc.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, usu)
}

func (uc *UsuariosController) getUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	usu := usuario{ID: id}
	if err := usu.getUsuario(uc.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, usu)
}

func (uc *UsuariosController) updateUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var usu usuario
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&usu); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	usu.ID = id

	if err := usu.updateUsuario(uc.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, usu)
}

func (uc *UsuariosController) deleteUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	usu := usuario{ID: id}
	if err := usu.deleteUsuario(uc.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
