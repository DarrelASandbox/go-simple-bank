package main

import (
	"database/sql"
	"log"

	"github.com/DarrelASandbox/go-simple-bank/db/util"

	"github.com/DarrelASandbox/go-simple-bank/api"
	db "github.com/DarrelASandbox/go-simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	// "." refers to current folder
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
