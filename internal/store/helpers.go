package store

import (
	"github.com/go-sql-driver/mysql"
)

func checkErrCode(err error, code uint16) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return me.Number == code
}
