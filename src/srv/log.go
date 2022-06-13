package srv

import (
	"github.com/gocql/gocql"

	"server/src"
	"server/src/log"
)

func LogAccess(code int, duration int, searchDuration int, error error, writeErr error, method string, uri string, encoding Encoding) {
	//language=SQL
	query := src.Session.Query(
		"INSERT INTO server.access (id, uri, code, duration, searchDuration, method, error, writeErr, encoding) VALUES (?,?,?,?,?,?,?,?,?)",
		gocql.TimeUUID(), uri, code, duration, searchDuration, method, (func() any {
			if error != nil {
				return error.Error()
			} else {
				return nil
			}
		})(), (func() any {
			if writeErr != nil {
				return writeErr.Error()
			} else {
				return nil
			}
		})(), encoding)
	err := query.Exec()
	if err != nil {
		log.Err(err, "Error inserting access into DB")
		log.Debug(query.Context())
	}
	log.Debug("LogAccess", uri, code, duration, searchDuration, method, error, writeErr, encoding)
}

/*
func LogAPIAccess(duration int, error error, request string) {
	//language=SQL
	query := src.Session.Query(
		"INSERT INTO server.apiaccess (id, duration, error, request) VALUES (?,?,?,?)",
		gocql.TimeUUID(), duration, (func() any {
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
	log.Debug("LogAPIAccess", duration, error, request)
}
*/
