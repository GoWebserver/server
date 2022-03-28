package src

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"

	"server/src/config"
	"server/src/log"
)

var cluster *gocql.ClusterConfig
var GQLSession gocqlx.Session
var Session *gocql.Session

// DBInit Create and open DB Connection
func DBInit() {
	cluster = gocql.NewCluster(config.GetConfig().Database.Hosts...)
	cluster.Port = int(config.GetConfig().Database.Port)
	cluster.Keyspace = config.GetConfig().Database.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.GetConfig().Database.User,
		Password: config.GetConfig().Database.Password,
	}
	var err error
	GQLSession, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Err(err, "Error creating connection")
		panic(err)
	}
	log.Log("GQLSession Connection established")
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Err(err, "Error creating connection")
		panic(err)
	}
	log.Log("Session Connection established")
}

func GetMimeTypes() map[string]string {
	now := time.Now()
	//language=SQL
	iter := Session.Query(
		"SELECT extension, mimetype FROM server.mime",
	).Iter()
	types := make(map[string]string, iter.NumRows())
	for {
		row := make(map[string]any)
		if !iter.MapScan(row) {
			break
		}
		types[fmt.Sprintf("%s", row["extension"])] = fmt.Sprintf("%s", row["mimetype"])
	}
	if err := iter.Close(); err != nil {
		log.Err(err, "Error loading Mime from DB")
		log.Debug(iter.Warnings())
	}
	log.Debug("Loaded Mime from DB in", time.Since(now))
	return types
}
