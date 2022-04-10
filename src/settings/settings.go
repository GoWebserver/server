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

	// min Size a file has to surpass,
	// in order to get compressed with Deflate
	//
	// default 1400
	DeflateCompressMinSize setting[uint64]

	// min Size a file has to surpass,
	// in order to get compressed with GZip
	//
	// default 1400
	GZipCompressMinSize setting[uint64]

	// min Size a file has to surpass,
	// in order to get compressed with Brotli
	//
	// default 1400
	BrotliCompressMinSize setting[uint64]

	// min Compression a file has to achieve,
	// in order to get compressed with Deflate
	//
	// default 0.2
	DeflateCompressMinCompression setting[float32]

	// min Compression a file has to achieve,
	// in order to get compressed with GZip
	//
	// default 0.2
	GZipCompressMinCompression setting[float32]

	// min Compression a file has to achieve,
	// in order to get compressed with Brotli
	//
	// default 0.2
	BrotliCompressMinCompression setting[float32]

	// Enable Deflate Compression
	//
	// default false
	EnableDeflateCompression setting[bool]

	// Enable GZip Compression
	//
	// default true
	EnableGZipCompression setting[bool]

	// Enable Brotli Compression
	//
	// default true
	EnableBrotliCompression setting[bool]
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
	sett.GZipCompressMinSize = setting[uint64]{
		Data:     1400,
		liveTime: onReload,
	}
	sett.DeflateCompressMinSize = setting[uint64]{
		Data:     1400,
		liveTime: onReload,
	}
	sett.BrotliCompressMinSize = setting[uint64]{
		Data:     1400,
		liveTime: onReload,
	}
	sett.GZipCompressMinCompression = setting[float32]{
		Data:     0.2,
		liveTime: onReload,
	}
	sett.DeflateCompressMinCompression = setting[float32]{
		Data:     0.2,
		liveTime: onReload,
	}
	sett.BrotliCompressMinCompression = setting[float32]{
		Data:     0.2,
		liveTime: onReload,
	}
	sett.EnableDeflateCompression = setting[bool]{
		Data:     false,
		liveTime: onReload,
	}
	sett.EnableGZipCompression = setting[bool]{
		Data:     true,
		liveTime: onReload,
	}
	sett.EnableBrotliCompression = setting[bool]{
		Data:     true,
		liveTime: onReload,
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
