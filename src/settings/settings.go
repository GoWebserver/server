package settings

import (
	"fmt"
	"time"

	"server/src/log"
)

type settings struct {
	// which site to serve if no path was specified
	// most likely be index.html
	//
	// default: "/index.html"
	DefaultSite setting[string]

	// list of all mimetypes for settings
	// the content type on response headers
	//
	// default empty
	Mimetypes setting[[]Mime]

	// setting to let the server only return off page
	//
	// default false
	ServerOff setting[bool]

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
	// value this setting contains
	data T

	// data that gets returned if no data could get loaded
	defaultData T

	// bool to indicate if this setting was already loaded
	loaded bool

	// lifetime of this setting
	liveTime LiveTime

	// data used for things like accesscount or last accesstime
	liveTimeData any

	// function to load Data from DB
	loadFunc func() error
}

type LiveTime uint8

const (
	// LoadEverytime loads setting from DB every time it gets accessed
	// never use this in anything serve Files related, only for rare requested
	// Settings where it doesn't matter if it takes ~ms to load
	// [Example]
	LoadEverytime LiveTime = iota

	// LoadAtStartup loads setting on first access or on forced settings reload
	// Use for Settings that significantly influence the server's behaviour
	// [Logging]
	LoadAtStartup

	// LoadAfterRequest reloads the setting async after every request
	// use for frequently changing and frequently requested Settings which
	// must be fast to access
	// [IndexPage, server shut down]
	LoadAfterRequest

	// LoadAfterXRequests reloads the setting async like LoadAfterRequest, but
	// only after certain count of requests were made.
	// use for frequently requested Settings which must be fast to access and
	// only sometimes change
	// [Mimetypes]
	LoadAfterXRequests

	// LoadAfterXTime reloads the setting op access if X time in ms has passed
	// since last access
	// use for settings which get rarely accessed, but if accessed ony times in
	// a short timespan
	// [Compression Info]
	LoadAfterXTime
)

type LoadAfterXRequestsData struct {
	XRequests     uint16
	countRequests uint16
}

type LoadAfterXTimeData struct {
	XTime      time.Duration
	lastAccess time.Time
}

var sett settings

func GetSettings() *settings {
	return &sett
}

func LoadDefaultSettings() {
	sett.DefaultSite = setting[string]{
		defaultData: "/index.html",
		liveTime:    LoadAfterRequest,
		loadFunc:    LoadDefaultSite,
	}
	sett.Mimetypes = setting[[]Mime]{
		defaultData: []Mime{},
		liveTime:    LoadAfterXRequests,
		liveTimeData: LoadAfterXRequestsData{
			XRequests: 1000,
		},
		loadFunc: LoadMimetypes,
	}
	sett.ServerOff = setting[bool]{
		defaultData: false,
		liveTime:    LoadAfterRequest,
		loadFunc:    LoadServerOff,
	}
	sett.DeflateCompressMinSize = setting[uint64]{
		defaultData: 1400,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadDeflateCompressMinSize,
	}
	sett.GZipCompressMinSize = setting[uint64]{
		defaultData: 1400,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadGZipCompressMinSize,
	}
	sett.BrotliCompressMinSize = setting[uint64]{
		defaultData: 1400,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadBrotliCompressMinSize,
	}
	sett.DeflateCompressMinCompression = setting[float32]{
		defaultData: 0.2,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadDeflateCompressMinCompression,
	}
	sett.GZipCompressMinCompression = setting[float32]{
		defaultData: 0.2,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadGZipCompressMinCompression,
	}
	sett.BrotliCompressMinCompression = setting[float32]{
		defaultData: 0.2,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadBrotliCompressMinCompression,
	}
	sett.EnableDeflateCompression = setting[bool]{
		defaultData: false,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadEnableDeflateCompression,
	}
	sett.EnableGZipCompression = setting[bool]{
		defaultData: true,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadEnableGZipCompression,
	}
	sett.EnableBrotliCompression = setting[bool]{
		defaultData: true,
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadEnableBrotliCompression,
	}
}

func (setting *setting[T]) Get() T {
	if !setting.loaded {
		err := setting.loadFunc()
		if err != nil {
			log.Err(err, fmt.Sprintf("Error loading initial Settings %#v using default data", setting))
			return setting.defaultData
		}
		setting.loaded = true
	}
	err := setting.load()
	if err != nil {
		log.Err(err, fmt.Sprintf("Error loading Settings %#v using default data", setting))
		return setting.defaultData
	}
	return setting.data
}

func (setting *setting[T]) load() (err error) {
	switch setting.liveTime {
	case LoadEverytime:
		err = setting.loadFunc()
	case LoadAfterRequest:
		go func() {
			err := setting.loadFunc()
			if err != nil {
				log.Err(err, fmt.Sprintf("Error loading Settings %#v async", setting))
			}
		}()
	case LoadAfterXRequests:
		data := setting.liveTimeData.(LoadAfterXRequestsData)
		data.countRequests++
		if data.countRequests >= data.XRequests {
			err = setting.loadFunc()
			if err == nil {
				data.countRequests = 0
			}
		}
		setting.liveTimeData = data
	case LoadAfterXTime:
		data := setting.liveTimeData.(LoadAfterXTimeData)
		if time.Since(data.lastAccess) > data.XTime {
			err = setting.loadFunc()
			if err == nil {
				data.lastAccess = time.Now()
			}
		}
		setting.liveTimeData = data
	}
	return
}
