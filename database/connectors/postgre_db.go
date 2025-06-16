package connectors

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PgSQLConn struct {
	DbHost     string
	DbPort     string
	DbDatabase string
	DbUsername string
	DbPassword string
}

func NewPgSQLConn(conn PgSQLConn) (*gorm.DB, error) {
	DbHost := conn.DbHost
	DbPort := conn.DbPort
	DbDatabase := conn.DbDatabase
	DbUsername := conn.DbUsername
	DbPassword := conn.DbPassword

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", DbUsername, DbPassword, DbHost, DbPort, DbDatabase)
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{Logger: newLogger})
}
