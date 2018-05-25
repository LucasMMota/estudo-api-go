package usuarios

import (
	"database/sql"
)

type usuario struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	//...
}

func getUsuarios(db *sql.DB, start, count int) ([]usuario, error) {
	rows, err := db.Query("SELECT * FROM usuarios LIMIT $1 OFFSET $2", count, start)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	usuarios := []usuario{}

	for rows.Next() {
		var usu usuario

		if err := rows.Scan(&usu.ID, &usu.Name, &usu.Email); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usu)
	}

	return usuarios, nil
}

func (usu *usuario) getUsuario(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM usuarios WHERE id=$1", usu.ID).Scan(&usu.ID, &usu.Name, &usu.Email)
}

func (usu *usuario) updateUsuario(db *sql.DB) error {
	_, err := db.Exec("UPDATE usuarios SET name=$1, email=$2 WHERE id=$3", usu.Name, usu.Email, usu.ID)

	return err
}

func (usu *usuario) deleteUsuario(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM usuarios WHERE id=$1", usu.ID)

	return err
}

func (usu *usuario) createUsuario(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO usuarios(name, email) VALUES($1,$2) RETURNING id", usu.Name, usu.Email).Scan(&usu.ID)

	if err != nil {
		return err
	}

	return nil
}
