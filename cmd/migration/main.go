package main

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/tanimutomo/sqlfile"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@go_db:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := sqlfile.New()
	err = s.File("deployment/migrations/init.sql")
	if err != nil {
		panic(err)
	}
	_, err = s.Exec(db)
	if err != nil {
		panic(err)
	}

	return
}
