package srv

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/andybalholm/brotli"

	"server/src/config"
	"server/src/log"
	"server/src/settings"
)

var root dir

func LoadSites() {
	log.Log("Loading Sites into Cache")
	start := time.Now()
	var size uint64
	var count counter
	wg := sync.WaitGroup{}
	root = loadDir(config.GetConfig().SitesDir, &size, &count, &wg)
	wg.Wait()
	log.Log(fmt.Sprintf("All files (%d) loaded in %s  Size:%dMB", count.count, time.Since(start), size/1048576))
	log.Debug(fmt.Sprintf("%d raw; %d deflate; %d gzip; %d br", count.rawcount, count.deflatecount, count.gzipcount, count.brcount))
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

func (data *data) getSmallest(encodings *map[Encoding]bool) (dat *[]byte, encoding Encoding) {
	dat = &data.raw
	var min = len(data.raw)
	if (*encodings)[Deflate] && data.deflate != nil && len(data.deflate) < min {
		min = len(data.deflate)
		dat = &data.deflate
		encoding = Deflate
	}
	if (*encodings)[GZip] && data.gzip != nil && len(data.gzip) < min {
		min = len(data.gzip)
		dat = &data.gzip
		encoding = GZip
	}
	if (*encodings)[Brotli] && data.br != nil && len(data.br) < min {
		min = len(data.br)
		dat = &data.br
		encoding = Brotli
	}
	return
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

type counter struct {
	count        uint32
	rawcount     uint32
	deflatecount uint32
	gzipcount    uint32
	brcount      uint32
}

func loadDir(path string, size *uint64, count *counter, wg *sync.WaitGroup) dir {
	siteCount, err := ioutil.ReadDir(path)
	if err != nil {
		log.Err(err, "Error reading directory", path)
		return dir{}
	}
	dir := dir{map[string]*file{}, map[string]dir{}}

	for _, site := range siteCount {
		site := site // prevents use uf loop variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			if site.IsDir() {
				dr := loadDir(fmt.Sprintf("%s/%s", path, site.Name()), size, count, wg)
				dir.dirs[site.Name()] = dr
				log.Debug(fmt.Sprintf("Loaded directory in cache %s/%s", path, site.Name()))
			} else {
				tmpSite, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, site.Name()))
				if err != nil {
					log.Err(err, fmt.Sprintf("Error loading site %s/%s", path, site.Name()))
				} else {
					file := createFile(tmpSite, site.Name(), count)
					*size += file.getSize()
					dir.files[site.Name()] = file
					log.Debug(fmt.Sprintf("Loaded site %s/%s  Size:%dMB", path, site.Name(), file.getSize()/1048576))
				}
			}
		}()
	}
	return dir
}

func createFile(raw []byte, name string, count *counter) *file {
	file := file{
		data: data{
			raw:     raw,
			deflate: nil,
			gzip:    nil,
			br:      nil,
		},
		mimetype: "",
	}
	count.rawcount++

	// ---------- deflate compress ----------
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
			log.Debug("compression to small for flate", float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100, "% ", name)
		}
		log.Debug("compressTime:", int(time.Since(now).Milliseconds()), "ms  deflate")
		count.deflatecount++
	}

	// ---------- gzip compress ----------
	if uint64(len(raw)) > settings.GetSettings().GZipCompressMinSize.Get() && settings.GetSettings().EnableGZipCompression.Get() {
		now := time.Now()
		var buf bytes.Buffer
		writer, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
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
			log.Debug("compression to small for gzip", float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100, "% ", name)
		}
		log.Debug("compressTime:", int(time.Since(now).Milliseconds()), "ms  gzip")
		count.gzipcount++
	}

	// ---------- br compress ----------
	if uint64(len(raw)) > settings.GetSettings().BrotliCompressMinSize.Get() && settings.GetSettings().EnableBrotliCompression.Get() {
		now := time.Now()
		var buf bytes.Buffer
		writer := brotli.NewWriterLevel(&buf, brotli.BestCompression)
		_, err := writer.Write(raw)
		if err != nil {
			log.Err(err, fmt.Sprintf("Error Brotling file %s", name))
		}
		if err := writer.Close(); err != nil {
			log.Err(err, fmt.Sprintf("Error Brotling file %s", name))
		}
		if float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100 > settings.GetSettings().BrotliCompressMinCompression.Get() {
			file.data.br = buf.Bytes()
		} else {
			log.Debug("compression to small for brotli", float32(len(file.data.raw)-buf.Len())/float32(len(file.data.raw))*100, "% ", name)
		}
		log.Debug("compressTime:", int(time.Since(now).Milliseconds()), "ms  brotli")
		count.brcount++
	}

	// ---------- mimetype ----------
	fileSplit := strings.Split(name, ".")
	filetype := fileSplit[len(fileSplit)-1]
	if val, exists := getMime(filetype); exists {
		file.mimetype = val
	} else {
		file.mimetype = ""
		log.Debug("unknown MimeType for extension", filetype)
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
	count.count++
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
