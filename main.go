package main

import (
	"BidServer/constants"
	"BidServer/infras/config"
	"BidServer/infras/connection"
	"BidServer/infras/logger"
	apis "BidServer/infras/router"
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/jmoiron/sqlx"
)

func initLastMigration(DbSource string) {
	migration, err := migrate.New(constants.MIGRATION_FOLDER_PATH, DbSource)
	if err != nil {
		log.Fatal("Cannot create migration: ", err)
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Cannot migrate db: ", err)
	}
	fmt.Println("Migrate succesfully")
}

func main() {
	//Init env config
	config, err := config.InitConfig(constants.ENV_FILE_PATH)
	if err != nil {
		log.Fatal("Internal server error: Cannot init config")
	}

	//Init logger
	_, err = logger.InitLogger(constants.LOG_FILE_PATH, config.Enviroment)
	if err != nil {
		log.Fatal("Internal server error: Cannot init logger")
	}

	//Init DB connection
	sqlContext, sqlxContext := connection.InitConnection(config.DbDriver, config.DbSource)
	defer func(sqlContext *sql.DB) {
		sqlContext.Close()
	}(sqlContext)

	defer func(sqlxContext *sqlx.DB) {
		sqlxContext.Close()
	}(sqlxContext)

	//DEV ONLY: Init last migration
	if config.Enviroment == constants.DEV_ENV_KEY {
		initLastMigration(config.DbSource)
	}

	//Init web apis
	server, err := apis.InitServer(sqlxContext)
	if err != nil {
		log.Fatal("Internal server error: Cannot start server")
	}

	//BEGIN THE APP - GLHF!
	server.Start()

}
