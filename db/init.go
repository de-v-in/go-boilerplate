package db

import (
	"fmt"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/phucvinh57/go-crud-example/configs"
)

func Init() *sql.DB {
	config := configs.GetConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.PgHost, config.PgPort, config.PgUser, config.PgPassword, config.PgDB)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connect to database successfully!")
	return db
}
