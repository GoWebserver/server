package srv

import (
	"regexp"
	"time"

	"server/src"
	"server/src/log"
)

var mimeTypes map[*regexp.Regexp]string

func LoadMime() {
	now := time.Now()
	types := src.GetMimeTypes()
	mimeTypes = make(map[*regexp.Regexp]string, len(types))
	for ext, typ := range types {
		mimeTypes[regexp.MustCompile(ext)] = typ
	}
	log.Log("Loaded Mime Types in", time.Since(now))
}

func getMime(filename string) (string, bool) {
	for r, s := range mimeTypes {
		if r.Match([]byte(filename)) {
			return s, true
		}
	}
	return "", false
}
