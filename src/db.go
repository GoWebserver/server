package src

import (
	"fmt"
	"strconv"
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

func GetMimeTypes() []Mime {
	now := time.Now()
	//language=SQL
	iter := Session.Query(
		"SELECT regex, type, \"index\" FROM server.mime",
	).Iter()
	types := make([]Mime, iter.NumRows())
	for {
		row := make(map[string]any)
		if !iter.MapScan(row) {
			break
		}
		index, _ := strconv.Atoi(fmt.Sprintf("%d", row["index"]))
		types[index] = Mime{
			Regex: fmt.Sprintf("%s", row["regex"]),
			Type:  fmt.Sprintf("%s", row["type"]),
		}
	}
	if err := iter.Close(); err != nil {
		log.Err(err, "Error loading Mime from DB")
		log.Debug(iter.Warnings())
	}
	log.Debug("Loaded Mime from DB in", time.Since(now))
	return types
}

type Mime struct {
	Regex string
	Type  string
}
