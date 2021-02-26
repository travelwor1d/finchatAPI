package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nkanders/finchat-api/internal/appconfig"

	// Importing mySQL driver
	_ "github.com/go-sql-driver/mysql"
)

func Connect(conf appconfig.MySQL) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.DB)
	return sqlx.Connect("mysql", dsn)
}
