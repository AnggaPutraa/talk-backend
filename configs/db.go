package configs

import (
	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func OpenConnection(source string) *sql.DB {
	connection, err := sql.Open("postgres", source)
	if err != nil {
		log.Fatal("Can't connect to db: ", err)
	}
	return connection
}
