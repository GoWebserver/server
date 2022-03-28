package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

// DB struct containing information about the Database connection
type DB struct {

	// Hosts of DB to connect to as an array.
	// Database to store logs, access logs, etc.
	//
	// default: []string{"no host provided"}
	Hosts []string `yaml:"Hosts"`

	// Port of DB to connect to.
	// Database to store logs, access logs, etc.
	//
	// default: 0
	Port uint16 `yaml:"Port"`

	// User of DB to connect to.
	// Database to store logs, access logs, etc.
	//
	// default: "no user provided"
	User string `yaml:"User"`

	// Password of User to connect to.
	// Database to store logs, access logs, etc.
	//
	// default: "no password provided"
	Password string `yaml:"Password"`

	// Keyspace of DB to use.
	// Database to store logs, access logs, etc.
	//
	// default: "no keyspace provided"
	Keyspace string `yaml:"Keyspace"`
}

type config struct {
	// PortHTTPS for the website must be between 0 and 65536
	// this comes from the Dockerfile and should
	// not get changed via the config file if used with Docker,
	// as it gets overridden with the environment variables
	//
	// default: 8443
	Port uint16 `yaml:"PortHTTPS" env:"Port"`

	// ApiPort used for the api must be between 0 and 65536
	// should be different from Port to avoid trying to serve
	// api by server
	//
	// default: 18266
	ApiPort uint16 `yaml:"ApiPort" env:"ApiPort"`

	// change to serve root for serving files
	// can be relative to the server main.go
	// or absolute
	//
	// only this directory is served, but no underlying directory
	// get served
	//
	// default: ./site
	SitesDir string `yaml:"SitesDir"`

	// removes Debug logs from console if set to true
	// disabling improves cache loading and serving speed
	//
	// default: false
	Debug bool `yaml:"Debug"`

	// Configuration for Database connection
	//
	// see DB
	Database DB `yaml:"Database"`
}

const (
	ConfigFile = "server.yml"
	CertsFile  = "certs/cert.pem"
	KeyFile    = "certs/key.key"
)

var conf config

func GetConfig() *config {
	return &conf
}

func LoadConfig() {
	defaultConfig()

	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		fmt.Printf("Error reading %s err:%s", ConfigFile, err)
		panic(err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		fmt.Printf("Error unmarshalling configs %s err:%s", ConfigFile, err)
		panic(err)
	}

	// load some values from env (for docker)
	loadEnv(&conf)
}

func defaultConfig() {
	conf.Port = 8443
	conf.ApiPort = 18266

	conf.SitesDir = "./site"

	conf.Debug = false

	conf.Database = DB{
		Hosts:    []string{"no host provided"},
		Port:     0,
		User:     "no user provided",
		Keyspace: "no keyspace provided",
		Password: "no database provided",
	}
}

func loadEnv(cfg *config) {
	confType := reflect.TypeOf(*cfg)
	out := reflect.New(reflect.Indirect(reflect.ValueOf(*cfg)).Type()).Elem()
	old := reflect.ValueOf(*cfg)
	for i := 0; i < confType.NumField(); i++ {
		outField := out.Field(i)
		confFieldType := confType.Field(i)
		oldField := old.Field(i)

		// copy old value
		outField.Set(oldField)

		// find env name
		tag, ok := confFieldType.Tag.Lookup("env")
		if !ok {
			continue
		}

		// fing env
		if env := os.Getenv(tag); env != "" {
			switch confFieldType.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				data, err := strconv.Atoi(env)
				if err == nil {
					outField.SetUint(uint64(data))
				}
			default:
				outField.SetString(env)
			}
		}
	}
	*cfg = out.Interface().(config)
}
