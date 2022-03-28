package srv

import (
	"fmt"
	"io/ioutil"
	"os"
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
	files map[string][]byte
	dirs  map[string]dir
}

func loadDir(name string, size *uint64, count *uint32) dir {
	siteCount, err := ioutil.ReadDir(name)
	if err != nil {
		log.Err(err, "Error reading directory", name)
		return dir{}
	}
	dir := dir{map[string][]byte{}, map[string]dir{}}

	for _, site := range siteCount {
		if info, _ := os.Stat(name + "/" + site.Name()); info.IsDir() {
			dr := loadDir(name+"/"+site.Name(), size, count)
			dir.dirs[site.Name()] = dr
			log.Debug("Loaded directory in cache", fmt.Sprintf("%s/%s", name, site.Name()))
		} else {
			tmpSite, err := ioutil.ReadFile(name + "/" + site.Name())
			if err != nil {
				log.Err(err, "Error loading site", fmt.Sprintf("%s/%s", name, site.Name()))
			} else {
				*size += uint64(len(tmpSite))
				*count++
				dir.files[site.Name()] = tmpSite
				log.Debug("Loaded site in cache", fmt.Sprintf("%s/%s", name, site.Name()))
			}
		}
	}
	return dir
}
