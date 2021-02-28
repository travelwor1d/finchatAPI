package store

import (
	"github.com/finchatapp/finchat-api/internal/appconfig"
	"github.com/jmoiron/sqlx"

	// Importing mySQL driver
	_ "github.com/go-sql-driver/mysql"
)

func Connect(conf appconfig.MySQL) (*sqlx.DB, error) {
	return sqlx.Connect("mysql", conf.ConnectionString)
}
