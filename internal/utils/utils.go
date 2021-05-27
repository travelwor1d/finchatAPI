package utils

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type Pair struct {
	Key, Value string
}

func ConvertViaJSON(from, to interface{}) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, to)
}

func SqlErrLogMsg(err error, query string, params []interface{}) string {
	msg := fmt.Sprintf("sql err: '%s', sql: '%s', params: %+v", err, query, params)
	// if os.Getenv("IS_LOCALHOST") == "true" {
	// 	defer panic(msg)
	// }
	debug.PrintStack()

	return msg
}

func DuplicateError(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	if me.Number == 1062 {
		return true
	}
	return strings.Contains(err.Error(), "duplicate")
}
func ConnectionResetByPeerError(err error) bool {
	return strings.Contains(err.Error(), "connection reset by peer")
}

func ColNamesWithPref(cols []string, pref string) []string {
	prefcols := make([]string, len(cols))
	copy(prefcols, cols)
	sort.Strings(prefcols)
	if pref == "" {
		return prefcols
	}

	for i := range prefcols {
		if !strings.Contains(prefcols[i], ".") {
			prefcols[i] = fmt.Sprintf("%s.%s", pref, prefcols[i])
		}
	}

	return prefcols
}
