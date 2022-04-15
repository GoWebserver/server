package srv

import (
	"time"

	"server/src/log"
	"server/src/settings"
)

func getMime(filename string) (string, bool) {
	start := time.Now()
	index := 0
	for _, s := range settings.GetSettings().Mimetypes.Get() {
		if s.Regex.Match([]byte(filename)) {
			log.Debug(float32(time.Since(start).Microseconds()))
			log.Debug(index)
			return s.Type, true
		}
		index++
	}
	log.Debug(float32(time.Since(start).Microseconds()))
	log.Debug("dead")
	return "", false
}
