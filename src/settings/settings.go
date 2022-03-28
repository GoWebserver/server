package settings

import (
	"regexp"
	"time"

	"server/src"
	"server/src/log"
)

type settings struct {
	// which site to serve if no path was specified
	// most likely be index.html
	//
	// default: "/index.html"
	DefaultSite setting[string]

	// at what rate to reload configs from
	// the Database in milliseconds
	//
	// default: 600000
	ReloadTime setting[uint32]

	// list of all mimetypes for settings
	// the content type on response headers
	//
	// default empty
	Mimetypes setting[[]Mime]
}

type setting[T any] struct {
	Data     T
	liveTime LiveTime
}

type LiveTime uint8

const (
	everytime LiveTime = iota
	atStartup
	onReload
)

var sett settings

func GetSettings() *settings {
	return &sett
}

func LoadDefaultSettings() {
	sett.DefaultSite = setting[string]{
		Data:     "/index.html",
		liveTime: onReload,
	}
	sett.ReloadTime = setting[uint32]{
		Data:     600000,
		liveTime: atStartup,
	}
	sett.Mimetypes = setting[[]Mime]{
		Data:     []Mime{},
		liveTime: atStartup,
	}
}

type Mime struct {
	Regex *regexp.Regexp
	Type  string
}

func LoadMime() {
	now := time.Now()
	types := src.GetMimeTypes()
	sett.Mimetypes.Data = make([]Mime, len(types))
	for ext, mime := range types {
		sett.Mimetypes.Data[ext] = Mime{regexp.MustCompile(mime.Regex), mime.Type}
	}
	log.Log("Loaded Mime Types in", time.Since(now))
}
