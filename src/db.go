package src

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"

	"server/src/config"
	"server/src/log"
)

var cluster *gocql.ClusterConfig
var GQLSession gocqlx.Session
var logsession *gocql.Session

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
	log.Log("Connection established")
	logsession, err = cluster.CreateSession()
	if err != nil {
		log.Err(err, "Error creating connection")
		panic(err)
	}
	log.Log("Connection established")
}
