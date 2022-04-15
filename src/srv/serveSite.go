package srv

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"server/src/log"
	"server/src/settings"
)

func getSite(path string, host string, availableEncodings *map[Encoding]bool) (*[]byte, int, error) {
	/*for _, forbidden := range util.GetConfig().Forbidden.Endpoints {
		if strings.HasPrefix(path, forbidden+"/") || path == forbidden {
			site, code := GetErrorSite(Forbidden, host, path)
			return &site, code, errors.New(path + " Forbidden by Endpoints " + forbidden)
		}
	}

	for _, forbidden := range util.GetConfig().Forbidden.Regex {
		match, err := regexp.Match(forbidden, []byte(path))
		if err != nil {
			// site, code := GetErrorSite(InternalServerError, host, path, fmt.Sprintf("Error checking forbidden regex"))
			return &site, code, err
		}
		if match {
			// site, code := GetErrorSite(Forbidden, host, path)
			return &site, code, errors.New(path + " Forbidden by regex " + forbidden)
		}
	}*/

	pathSplit := strings.Split(path, "/")[1:]

	depth := len(pathSplit)

	dir := root
	for i := 0; i < depth-1; i++ {
		dir = dir.dirs[pathSplit[i]]
		if dir.files == nil {
			break
		}
	}
	site, ok := dir.files[pathSplit[depth-1]]
	if !ok {
		if _, ok := dir.dirs[pathSplit[depth-1]]; ok {
			// site, code := GetErrorSite(NotFound, host, path, fmt.Sprintf("%s is no file, but a directory", pathSplit[depth-1]))
			return nil, 404, errors.New(fmt.Sprintf("no site data for: %v", pathSplit))
		}
		// site, code := GetErrorSite(NotFound, host, path)
		return nil, 404, errors.New(fmt.Sprintf("no site data for: %s", pathSplit))
	}
	return site.getSmallest(availableEncodings), 200, nil
}

// CreateServe
//
// Registers a handle for '/' to serve the DefaultSite
func CreateServe() http.HandlerFunc {
	fun := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if r.URL.Path == "/" {
			r.URL.Path = settings.GetSettings().DefaultSite.Get()
		}
		encodings := strings.Split(r.Header.Get("Accept-Encoding"), ",")
		availableEncodings := map[Encoding]bool{
			Deflate: false,
			GZip:    false,
			Brotli:  false,
		}
		for _, encoding := range encodings {
			availableEncodings[Encoding(strings.TrimSpace(encoding))] = true
		}

		msg, code, err := getSite(r.URL.Path, r.Host, &availableEncodings)

		searchTime := time.Now()

		if err != nil {
			log.Err(err, fmt.Sprintf("Error getting site %s", r.URL.Path))
			w.WriteHeader(code)
		} else {
			fileSplit := strings.Split(r.URL.Path[1:], ".")
			filetype := fileSplit[len(fileSplit)-1]
			if val, exists := getMime(filetype); exists {
				w.Header().Set("Content-Type", val)
			} else {
				log.Log("unknown MimeType", filetype)
			}
		}

		_, er := w.Write(*msg)
		if er != nil {
			log.Err(er, "Error writing response:")
		}
		go LogAccess(code, int(time.Since(start).Microseconds()), int(searchTime.Sub(start).Microseconds()), err, er, r.Method, r.URL.Path)
	}

	return fun
}

type Encoding string

const (
	GZip    Encoding = "gzip"
	Deflate Encoding = "deflate"
	Brotli  Encoding = "br"
)
