package settings

import (
	"fmt"
	"sync"
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

	// List of regexes or strings that prevent
	// access to certain routs
	//
	// default []
	Forbidden setting[[]Forbidden]

	// Maximum length of request URIs
	//
	// default: 1000
	MaxURILength setting[uint16]
}

type setting[T any] struct {
	// value this setting contains
	data T

	// data that gets returned if no data could get loaded
	defaultData T

	// bool to indicate if this setting was already loaded
	loaded bool

	// Mutex to indicate if this setting is currently loaded (prevent multiple loadings)
	loading sync.RWMutex

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
	// Settings where it doesn't matter if it takes ~ms to load and are not
	// allowed to be out of sync
	// []
	LoadEverytime LiveTime = iota

	// LoadAsyncAfterEveryRequest reloads the setting async after every request
	// use for frequently changing and frequently requested Settings which
	// must be fast to access
	// [IndexPage, server shut down]
	LoadAsyncAfterEveryRequest

	// LoadAsyncAfterXRequestsAfterRequest reloads the setting async like LoadAsyncAfterEveryRequest, but
	// only after certain count of requests were made.
	// use for frequently requested Settings which must be fast to access and
	// only sometimes change
	// []
	LoadAsyncAfterXRequestsAfterRequest

	// LoadAfterXTime reloads the setting on access if X time in ms has passed
	// since last access
	// use for settings which get rarely accessed, but if accessed many times in
	// a short timespan
	// [Compression Info, Mimetypes]
	LoadAfterXTime

	// LoadAfterXTimeAfterAccess reloads the setting after access if X time in ms has passed
	// since last access
	// use for settings which get frequently accessed, doesn't change often, is allowed
	// to be out of date on first request after some time and must be fast to access
	// [Forbidden]
	LoadAfterXTimeAfterAccess
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
		liveTime:    LoadAsyncAfterEveryRequest,
		loadFunc:    LoadDefaultSite,
	}
	sett.Mimetypes = setting[[]Mime]{
		defaultData: []Mime{},
		liveTime:    LoadAfterXTime,
		liveTimeData: LoadAfterXTimeData{
			XTime: 60 * time.Second,
		},
		loadFunc: LoadMimetypes,
	}
	sett.ServerOff = setting[bool]{
		defaultData: false,
		liveTime:    LoadAsyncAfterEveryRequest,
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
	sett.MaxURILength = setting[uint16]{
		defaultData: 1000,
		liveTime:    LoadAsyncAfterXRequestsAfterRequest,
		liveTimeData: LoadAfterXRequestsData{
			XRequests: 100,
		},
		loadFunc: LoadMaxURILength,
	}
	sett.Forbidden = setting[[]Forbidden]{
		defaultData: []Forbidden{},
		liveTime:    LoadAfterXTimeAfterAccess,
		liveTimeData: LoadAfterXTimeData{
			XTime: 30 * time.Second,
		},
		loadFunc: LoadForbidden,
	}
}

func (setting *setting[T]) Get() T {
	// block to check if loaded
	setting.loading.RLock()

	// only check if loaded after acquiring lock (so if setting was loaded in meantime this doesn't get called again
	if !setting.loaded {

		// release read to modify setting
		setting.loading.RUnlock()

		data, err := func() (T, error) {
			// lock the setting for writing to load
			setting.loading.Lock()
			defer setting.loading.Unlock()

			// check if setting was already loaded
			if !setting.loaded {
				err := setting.loadFunc()
				if err != nil {
					log.Err(err, fmt.Sprintf("Error loading initial Settings %#v using default data", setting))
					return setting.defaultData, nil
				}
				setting.loaded = true

				// update data
				switch setting.liveTime {
				case LoadAsyncAfterXRequestsAfterRequest:
					data := setting.liveTimeData.(LoadAfterXRequestsData)
					data.countRequests = 0
					setting.liveTimeData = data
				case LoadAfterXTime, LoadAfterXTimeAfterAccess:
					data := setting.liveTimeData.(LoadAfterXTimeData)
					data.lastAccess = time.Now()
					setting.liveTimeData = data
				}
				return setting.data, nil
			}
			return setting.data, fmt.Errorf("allready loaded")
		}()
		if err == nil {
			// return loaded data
			return data
		}
		log.Debug("setting already loaded")
		// continue with normal function if was already loaded
	} else {
		// release read after loaded check
		setting.loading.RUnlock()
	}

	err := setting.load()
	if err != nil {
		log.Err(err, fmt.Sprintf("Error loading Settings %#v using default data", setting))
		return setting.defaultData
	}

	// get lock to read data
	setting.loading.RLock()
	defer setting.loading.RUnlock()
	return setting.data
}

func (setting *setting[T]) load() (err error) {
	switch setting.liveTime {
	case LoadEverytime:
		setting.loading.Lock()
		defer setting.loading.Unlock()
		err = setting.loadFunc()
	case LoadAsyncAfterEveryRequest:
		go func() {
			setting.loading.Lock()
			defer setting.loading.Unlock()
			err := setting.loadFunc()
			if err != nil {
				log.Err(err, fmt.Sprintf("Error loading Setting %#v async", setting))
			}
		}()
	case LoadAsyncAfterXRequestsAfterRequest:
		data := setting.liveTimeData.(LoadAfterXRequestsData)
		data.countRequests++
		if data.countRequests >= data.XRequests {
			go func() {
				setting.loading.Lock()
				defer setting.loading.Unlock()
				err := setting.loadFunc()
				if err != nil {
					log.Err(err, fmt.Sprintf("Error loading Setting %#v async", setting))
				} else {
					data.countRequests = 0
				}
			}()
		}
		setting.liveTimeData = data
	case LoadAfterXTime:
		data := setting.liveTimeData.(LoadAfterXTimeData)
		if time.Since(data.lastAccess) > data.XTime {
			setting.loading.Lock()
			defer setting.loading.Unlock()
			err = setting.loadFunc()
			if err == nil {
				data.lastAccess = time.Now()
			}
		}
		setting.liveTimeData = data
	case LoadAfterXTimeAfterAccess:
		data := setting.liveTimeData.(LoadAfterXTimeData)
		if time.Since(data.lastAccess) > data.XTime {
			go func() {
				setting.loading.Lock()
				defer setting.loading.Unlock()
				err = setting.loadFunc()
				if err == nil {
					data.lastAccess = time.Now()
				}
			}()
		}
		setting.liveTimeData = data
	}
	return
}
