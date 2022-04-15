package srv

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"server/src/config"
	"server/src/log"
	"server/src/settings"
)

var root dir

func LoadSites() {
	log.Log("Loading Sites into Cache")
	start := time.Now()
	var size uint64
	var count uint32
	root = loadDir(config.GetConfig().SitesDir, &size, &count)
	log.Log(fmt.Sprintf("All files (%d) loaded in %s  Size:%dMB", count, time.Since(start), size/1048576))
	runtime.GC()
}

type dir struct {
	files map[string]*file
	dirs  map[string]dir
}

type data struct {
	raw     []byte
	gzip    []byte
	deflate []byte
	br      []byte
}

type file struct {
	data     data
	mimetype string
}

func (data *data) getSmallest(encodings *map[Encoding]bool) *[]byte {
	// order matters
	if (*encodings)[Deflate] && data.deflate != nil {
		return &data.deflate
	}
	if (*encodings)[GZip] && data.gzip != nil {
		return &data.gzip
	}
	if (*encodings)[Brotli] && data.br != nil {
		return &data.br
	}

	return &data.raw
}

func (file *file) getSize() (size uint64) {
	size = uint64(len(file.data.raw))
	if file.data.deflate != nil {
		size += uint64(len(file.data.deflate))
	}
	if file.data.gzip != nil {
		size += uint64(len(file.data.gzip))
	}
	if file.data.br != nil {
		size += uint64(len(file.data.br))
	}
	return
}

func loadDir(path string, size *uint64, count *uint32) dir {
	siteCount, err := ioutil.ReadDir(path)
	if err != nil {
		log.Err(err, "Error reading directory", path)
		return dir{}
	}
	dir := dir{map[string]*file{}, map[string]dir{}}

	for _, site := range siteCount {
		if site.IsDir() {
			dr := loadDir(fmt.Sprintf("%s/%s", path, site.Name()), size, count)
			dir.dirs[site.Name()] = dr
			log.Debug(fmt.Sprintf("Loaded directory in cache %s/%s", path, site.Name()))
		} else {
			tmpSite, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, site.Name()))
			if err != nil {
				log.Err(err, fmt.Sprintf("Error loading site %s/%s", path, site.Name()))
			} else {
				file := createFile(tmpSite, site.Name())
				*size += file.getSize()
				*count++
				dir.files[site.Name()] = file
				log.Debug(fmt.Sprintf("Loaded site in cache %s/%s  Size:%dMB", path, site.Name(), file.getSize()/1048576))
			}
		}
	}
	return dir
}

func createFile(raw []byte, name string) *file {
	file := file{
		data: data{
			raw:     raw,
			deflate: nil,
			gzip:    nil,
			br:      nil,
		},
		mimetype: "",
	}

	// ---------- flate compress ----------
	if uint64(len(raw)) > settings.GetSettings().DeflateCompressMinSize.Get() && settings.GetSettings().EnableDeflateCompression.Get() {
		now := time.Now()
		var buf bytes.Buffer
		writer, err := flate.NewWriter(&buf, flate.BestCompression)
		if err != nil {
			log.Err(err, fmt.Sprintf("Error Flating file %s", name))
		}
		_, err = writer.Write(raw)
		if err != nil {
			log.Err(err, fmt.Sprintf("Error Flating file %s", name))
		}
		if err := writer.Close(); err != nil {
			log.Err(err, fmt.Sprintf("Error Flating file %s", name))
		}
		if float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100 > settings.GetSettings().DeflateCompressMinCompression.Get() {
			file.data.deflate = buf.Bytes()
		} else {
			log.Debug("compression to small for flate", float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100, name)
		}
		log.Debug("compressTime:", int(time.Since(now).Milliseconds()), "ms  flate")
	}

	// ---------- gzip compress ----------
	if uint64(len(raw)) > settings.GetSettings().GZipCompressMinSize.Get() && settings.GetSettings().EnableGZipCompression.Get() {
		now := time.Now()
		var buf bytes.Buffer
		writer, err := gzip.NewWriterLevel(&buf, flate.BestCompression)
		if err != nil {
			log.Err(err, fmt.Sprintf("Error GZipping file %s", name))
		}
		_, err = writer.Write(raw)
		if err != nil {
			log.Err(err, fmt.Sprintf("Error GZipping file %s", name))
		}
		if err := writer.Close(); err != nil {
			log.Err(err, fmt.Sprintf("Error GZipping file %s", name))
		}
		if float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100 > settings.GetSettings().GZipCompressMinCompression.Get() {
			file.data.gzip = buf.Bytes()
		} else {
			log.Debug("compression to small for gzip", float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100, name)
		}
		log.Debug("compressTime:", int(time.Since(now).Milliseconds()), "ms  gzip")
	}

	// ---------- br compress ----------
	//

	// ---------- mimetype ----------
	fileSplit := strings.Split(name, ".")
	filetype := fileSplit[len(fileSplit)-1]
	if val, exists := getMime(filetype); exists {
		file.mimetype = val
	} else {
		file.mimetype = ""
		log.Debug("unknown MimeType for extension"+
			"", filetype)
	}

	// ---------- log ----------
	log.Debug(fmt.Sprintf(
		"Loaded file %s with sizes: {raw:%dMB, flate:%s, gzip:%s, brotli:%s} mimetype:%s", name,
		len(file.data.raw)/1048576,
		(func() string {
			if file.data.deflate != nil {
				return fmt.Sprintf("%dMB %.2f%%compression", len(file.data.deflate)/1048576, float32(len(file.data.raw)-len(file.data.deflate))/float32(len(file.data.raw))*100)
			} else {
				return "no Compression"
			}
		})(),
		(func() string {
			if file.data.gzip != nil {
				return fmt.Sprintf("%dMB %.2f%%compression", len(file.data.gzip)/1048576, float32(len(file.data.raw)-len(file.data.gzip))/float32(len(file.data.raw))*100)
			} else {
				return "no Compression"
			}
		})(), (func() string {
			if file.data.br != nil {
				return fmt.Sprintf("%dMB %.2f%%compression", len(file.data.br)/1048576, float32(len(file.data.raw)-len(file.data.br))/float32(len(file.data.raw))*100)
			} else {
				return "no Compression"
			}
		})(), file.mimetype,
	))
	return &file
}

func getMime(filename string) (string, bool) {
	index := 0
	for _, s := range settings.GetSettings().Mimetypes.Get() {
		if s.Regex.Match([]byte(filename)) {
			return s.Type, true
		}
		index++
	}
	return "", false
}
