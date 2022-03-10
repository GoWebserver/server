package src

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Logging struct containing information about logging
type Logging struct {

	// adds a LogGroup to the log ( |CONFIG ) and adds a prefix to indicate
	// the type of log (> for Normal, * for Debug, ! for Error) (<-- default)
	//
	// default: false
	LogPrefix bool `yaml:"LogPrefix"`

	// stretches the prefix with LogGroup and > / * / ! to certain size
	//
	// can be ignored if LogPrefix is set to false
	//
	// default: 9 (fits the longest group)
	LogFile bool `yaml:"LogFile"`

	// adds the file to the log where the Log method was called
	// should be activated for debug purposes
	//
	// default: false
	StretchPrefix int8 `yaml:"StretchPrefix"`

	// stretches the filename:line number to certain size
	//
	// can be ignored if LogFile is set to false
	//
	// default: 16
	StretchFile int8 `yaml:"StretchFile"`
}

// DB struct containing information about the Database connection
type DB struct {

	// host of DB to connect to.
	// Database to store logs, access logs, etc.
	//
	// default: "no host provided"
	Host string `yaml:"Host"`

	// port of DB to connect to.
	// Database to store logs, access logs, etc.
	//
	// default: "no port provided"
	Port uint16 `yaml:"Port"`

	// user of DB to connect to.
	// Database to store logs, access logs, etc.
	//
	// default: "no user provided"
	User string `yaml:"User"`

	// password of User to connect to.
	// Database to store logs, access logs, etc.
	//
	// default: "no password provided"
	Password string `yaml:"Password"`

	// database of DB to use.
	// Database to store logs, access logs, etc
	//
	// default: "no database provided"
	Database string `yaml:"Database"`
}

type config struct {

	// Port for the website must be between 0 and 65536
	// this comes from the Dockerfile and should
	// not get changed via the config file if used with Docker
	//
	// default: 80
	PortHTTP uint16

	// Port for the website must be between 0 and 65536
	// this comes from the Dockerfile and should
	// not get changed via the config file if used with Docker
	//
	// default: 443
	PortHTTPS uint16

	// Port used for the api must be between 0 and 65536
	// should be different from Port to avoid trying to serve
	// api by server
	//
	// default: 18266
	ApiPort uint16

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
	// can improve cache loading speed
	//
	// default: false
	Debug bool `yaml:"Debug"`

	// set of configurations fo logging
	//
	// see Logging
	Logging Logging `yaml:"Logging"`

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
	if os.Getenv("PortHTTP") != "" {
		port, err := strconv.Atoi(os.Getenv("PortHTTP"))
		if err != nil {
			conf.PortHTTP = uint16(port)
		}
	}
	if os.Getenv("PortHTTPS") != "" {
		port, err := strconv.Atoi(os.Getenv("PortHTTPS"))
		if err != nil {
			conf.PortHTTPS = uint16(port)
		}
	}
	if os.Getenv("APIPORT") != "" {
		apiPort, err := strconv.Atoi(os.Getenv("APIPORT"))
		if err != nil {
			conf.ApiPort = uint16(apiPort)
		}
	}
}

func defaultConfig() {
	conf.PortHTTP = 8080
	conf.PortHTTPS = 8443
	conf.ApiPort = 18266

	conf.SitesDir = "./site"
	conf.Logging = Logging{
		LogPrefix:     false,
		LogFile:       false,
		StretchPrefix: 9,
		StretchFile:   16,
	}

	conf.Debug = false

	conf.Database = DB{
		Host:     "no host provided",
		Port:     0,
		User:     "no user provided",
		Database: "no password provided",
		Password: "no database provided",
	}
}
