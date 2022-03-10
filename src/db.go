package src

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"

	"server/src/log"
)

var cluster *gocql.ClusterConfig
var GQLSession gocqlx.Session
var logsession *gocql.Session

// DBInit Create and open DB Connection
func DBInit() {
	cluster = gocql.NewCluster(GetConfig().Database.Host)
	cluster.Port = int(GetConfig().Database.Port)
	cluster.Keyspace = GetConfig().Database.Database
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: GetConfig().Database.User,
		Password: GetConfig().Database.Password,
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
