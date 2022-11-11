package bootstrap

import (
	_ "github.com/go-sql-driver/mysql"
	"skframe/pkg/database"
)

func SetupDB() {
	database.ConnectDB()
}

func DestructDB() {
	database.Destruct()
}
