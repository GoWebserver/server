package srv

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"

	"server/src/config"
	"server/src/log"
	"server/src/settings"
)

var root dir

func LoadSites() error {
	log.Log("Loading Sites into Cache")
	start := time.Now()
	var size uint64
	var count uint32
	root = loadDir(config.GetConfig().SitesDir, &size, &count)
	log.Log(fmt.Sprintf("All files (%d) loaded in %s  Size:%dMB", count, time.Since(start), size/1048576))
	runtime.GC()

	return nil
}

type dir struct {
	files map[string]*file
	dirs  map[string]dir
}

type file struct {
	raw     []byte
	gzip    []byte
	deflate []byte
	br      []byte
}

func (file *file) getSmallest(encodings *map[Encoding]bool) *[]byte {
	// order matters
	if (*encodings)[Deflate] && file.deflate != nil {
		return &file.deflate
	}
	if (*encodings)[GZip] && file.gzip != nil {
		return &file.gzip
	}
	if (*encodings)[Brotli] && file.br != nil {
		return &file.br
	}

	return &file.raw
}

func (file *file) getSize() (size uint64) {
	size = uint64(len(file.raw))
	if file.deflate != nil {
		size += uint64(len(file.deflate))
	}
	if file.gzip != nil {
		size += uint64(len(file.gzip))
	}
	if file.br != nil {
		size += uint64(len(file.br))
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
	f := file{
		raw:     raw,
		deflate: nil,
		gzip:    nil,
		br:      nil,
	}

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
		if float32(len(f.raw)-buf.Len())/float32(len(f.raw))*100 > settings.GetSettings().DeflateCompressMinCompression.Get() {
			f.deflate = buf.Bytes()
		} else {
			log.Debug("compression to small for flate", float32(len(f.raw)-buf.Len())/float32(len(f.raw))*100, name)
		}
		log.Debug("compressTime:", int(time.Since(now).Milliseconds()), "ms  flate")
	}

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
		if float32(len(f.raw)-buf.Len())/float32(len(f.raw))*100 > settings.GetSettings().GZipCompressMinCompression.Get() {
			f.gzip = buf.Bytes()
		} else {
			log.Debug("compression to small for gzip", float32(len(f.raw)-buf.Len())/float32(len(f.raw))*100, name)
		}
		log.Debug("compressTime:", int(time.Since(now).Milliseconds()), "ms  gzip")
	}

	log.Debug(fmt.Sprintf(
		"Created file %s with sizes: {raw:%dMB, flate:%s, gzip:%s, brotli:%s}", name,
		len(f.raw)/1048576,
		(func() string {
			if f.deflate != nil {
				return fmt.Sprintf("%dMB %.2f%%compression", len(f.deflate)/1048576, float32(len(f.raw)-len(f.deflate))/float32(len(f.raw))*100)
			} else {
				return "no Compression"
			}
		})(),
		(func() string {
			if f.gzip != nil {
				return fmt.Sprintf("%dMB %.2f%%compression", len(f.gzip)/1048576, float32(len(f.raw)-len(f.gzip))/float32(len(f.raw))*100)
			} else {
				return "no Compression"
			}
		})(), (func() string {
			if f.br != nil {
				return fmt.Sprintf("%dMB %.2f%%compression", len(f.br)/1048576, float32(len(f.raw)-len(f.br))/float32(len(f.raw))*100)
			} else {
				return "no Compression"
			}
		})(),
	))
	return &f
}
