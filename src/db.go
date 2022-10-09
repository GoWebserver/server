package src

import (
	"github.com/gocql/gocql"

	"server/src/config"
	"server/src/log"
)

var cluster *gocql.ClusterConfig

// var GQLSession gocqlx.Session

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
	// GQLSession, err = gocqlx.WrapSession(cluster.CreateSession())
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
