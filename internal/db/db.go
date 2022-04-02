package db

import (
	"context"
	"github.com/ashirko/imlogs/internal/utils"
	immudb "github.com/codenotary/immudb/pkg/client"
	"log"
	"time"
)

const (
	DB_HOST             = "localhost"
	DB_PORT             = 3322
	DB_USER             = "immudb"
	DB_PASSWORD         = "immudb"
	DATABASE            = "defaultdb"
	MAX_DB_CONN_RETRIES = 5
)

type Database struct {
	Client immudb.ImmuClient
}

var db Database

func Initialize() (Database, error) {

	db = newDatabase()

	err := db.openSession()
	if err != nil {
		return db, err
	}

	err = db.createLogsTable()
	return db, err
}

func Close(db Database) {
	log.Println("Close database...")
	err := db.Client.CloseSession(context.TODO())
	if err != nil {
		log.Println("Error closing database session: ", err)
	}
}

func GetDB() Database {
	// TODO implement reconnection in case of broken connection
	return db
}

func newDatabase() Database {
	db := Database{}
	db.Client = immudb.NewClient()
	opts := getDatabaseOpts()
	db.Client = db.Client.WithOptions(opts)
	return db
}

func getDatabaseOpts() *immudb.Options {
	address := utils.GetEnvStr("DB_HOST", DB_HOST)
	port := utils.GetEnvInt("DB_PORT", DB_PORT)
	return immudb.DefaultOptions().WithAddress(address).WithPort(port)
}

func (db Database) openSession() (err error) {
	maxRetries := utils.GetEnvInt("MAX_DB_CONN_RETRIES", MAX_DB_CONN_RETRIES)
	user := utils.GetEnvStr("DB_USER", DB_USER)
	password := utils.GetEnvStr("DB_PASSWORD", DB_PASSWORD)
	database := utils.GetEnvStr("DATABASE", DATABASE)

	for i := 0; i < maxRetries; i++ {
		err = db.Client.OpenSession(context.TODO(), []byte(user), []byte(password), database)
		if err == nil {
			return
		}
		log.Println("Opening session to the database failed... Reason: ", err)
		time.Sleep(1 * time.Second)
	}
	return
}
