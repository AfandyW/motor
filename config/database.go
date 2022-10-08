package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDatabase(configDB *DB) *sql.DB {
	conn := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable",
		configDB.Host,
		configDB.Username,
		configDB.Password,
		configDB.BaseUrl,
		configDB.Database,
	)

	db, err := sql.Open(configDB.Host, conn)

	if err != nil {
		fmt.Println("connect database return error")
		panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		fmt.Println("ping database return error")
		panic(err.Error())
	}

	return db
}
