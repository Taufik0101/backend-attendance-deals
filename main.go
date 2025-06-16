package main

import (
	"backend-attendance-deals/config"
	"backend-attendance-deals/database/connectors"
	"backend-attendance-deals/routes"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
)

var pgsqldb *gorm.DB

const defaultPort = "8080"

func init() {
	var err error

	_ = godotenv.Load()

	err = config.CheckEnv()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	pgsqldb = initPG()
	route := routes.NewRoute(pgsqldb)
	app := route.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	errLis := app.Listen(":" + port)
	if errLis != nil {
		return
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
}

func initPG() *gorm.DB {
	const INIT_STEP = 0
	const APP_SCHEMA_VERSION = 8

	pgsqlConn := connectors.PgSQLConn{
		DbHost:     config.GetEnv("DB_PGSQL_HOST", ""),
		DbPort:     config.GetEnv("DB_PGSQL_PORT", ""),
		DbDatabase: config.GetEnv("DB_PGSQL_DATABASE", ""),
		DbUsername: config.GetEnv("DB_PGSQL_USERNAME", ""),
		DbPassword: config.GetEnv("DB_PGSQL_PASSWORD", ""),
	}

	dbConn, err := connectors.NewPgSQLConn(pgsqlConn)
	if err != nil {
		fmt.Println(err)
	}

	if dbConn == nil {
		fmt.Println("failed connect to pgsqldb database")
	}

	db, err := dbConn.DB()
	if err != nil {
		fmt.Println("failed connect to database")
	}
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(100)
	//database pooling

	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	_, err = migrate.NewWithDatabaseInstance(
		"file://./database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		fmt.Println("failed migration 1 : ", err.Error())
	}

	return dbConn
}
