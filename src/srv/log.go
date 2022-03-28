package srv

import (
	"github.com/gocql/gocql"

	"server/src"
	"server/src/log"
)

func LogAccess(code int, duration int, searchDuration int, error error, writeErr error, method string, uri string) {
	//language=SQL
	query := src.Session.Query(
		"INSERT INTO server.access (id, uri, code, duration, searchDuration, method, error, writeErr) VALUES (?,?,?,?,?,?,?,?,?)",
		gocql.TimeUUID(), uri, code, duration, searchDuration, method, (func() interface{} {
			if error != nil {
				return error.Error()
			} else {
				return nil
			}
		})(), (func() interface{} {
			if writeErr != nil {
				return writeErr.Error()
			} else {
				return nil
			}
		})())
	err := query.Exec()
	if err != nil {
		log.Err(err, "Error inserting access into DB")
		log.Debug(query.Context())
	}
}

func LogAPIAccess(duration int, error error, request string) {
	//language=SQL
	query := src.Session.Query(
		"INSERT INTO server.apiaccess (id, duration, error, request) VALUES (?,?,?,?)",
		gocql.TimeUUID(), duration, (func() interface{} {
			if error != nil {
				return error.Error()
			} else {
				return nil
			}
		})(), request)
	err := query.Exec()
	if err != nil {
		log.Err(err, "Error inserting accessapi into DB")
		log.Debug(query.Context())
	}
}
