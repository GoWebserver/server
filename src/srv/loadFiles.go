package srv

import (
	"fmt"
	"io/ioutil"
	"time"

	"server/src/config"
	"server/src/log"
)

var root dir

func LoadSites() error {
	log.Log("Loading Sites into Cache")
	start := time.Now()
	var size uint64
	var count uint32
	root = loadDir(config.GetConfig().SitesDir, &size, &count)
	log.Log(fmt.Sprintf("All files (%d) loaded in %s  Size:%dMB", count, time.Since(start), size/1048576))

	return nil
}

type dir struct {
	files map[string]file
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
	if (*encodings)[Brotli] && file.br != nil {
		return &file.br
	}
	if (*encodings)[GZip] && file.gzip != nil {
		return &file.gzip
	}
	if (*encodings)[Deflate] && file.deflate != nil {
		return &file.deflate
	}

	return &file.raw
}

func loadDir(path string, size *uint64, count *uint32) dir {
	siteCount, err := ioutil.ReadDir(path)
	if err != nil {
		log.Err(err, "Error reading directory", path)
		return dir{}
	}
	dir := dir{map[string]file{}, map[string]dir{}}

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
				*size += uint64(site.Size())
				*count++
				dir.files[site.Name()] = createFile(&tmpSite)
				log.Debug(fmt.Sprintf("Loaded site in cache %s/%s  Size:%dMB", path, site.Name(), site.Size()/1048576))
			}
		}
	}
	return dir
}

func createFile(raw *[]byte) file {
	log.Debug(len(*raw))
	return file{
		raw:     *raw,
		gzip:    nil,
		deflate: nil,
		br:      nil,
	}
}
