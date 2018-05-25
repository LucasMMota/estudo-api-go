package main

import (
	"log"
	"net/http"

	"github.com/TecnoFIS/back-end-APIs/usuarios"
)

func main() {
	usuariosController := usuarios.UsuariosController{}
	usuariosController.Initialize(
		"estudo2", // os.Getenv("APP_DB_USERNAME"),
		"estudo2", // os.Getenv("APP_DB_PASSWORD"),
		"estudo2", // os.Getenv("APP_DB_NAME"),
	)

	log.Fatal(http.ListenAndServe(":8080", usuariosController.Router))
}
