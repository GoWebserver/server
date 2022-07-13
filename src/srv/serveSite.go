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

func getSite(request *http.Request, availableEncodings *map[Encoding]bool) (*[]byte, Encoding, int, string, error) {
	if request.Method != "GET" {
		data, code := GetErrorSite(http.StatusMethodNotAllowed, request.Host, request.URL.Path, "")
		return data, "", code, "text/html", errors.New(fmt.Sprintf("not get method (%v)", request.Method))
	}

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

	pathSplit := strings.Split(request.URL.Path, "/")[1:]

	depth := len(pathSplit)

	dir := root
	for i := 0; i < depth-1; i++ {
		dir = dir.dirs[pathSplit[i]]
		if dir.files == nil {
			break
		}
	}
	file, ok := dir.files[pathSplit[depth-1]]
	if !ok {
		if _, ok := dir.dirs[pathSplit[depth-1]]; ok {
			data, code := GetErrorSite(http.StatusNotFound, request.Host, request.URL.Path, fmt.Sprintf("%s is no file, but a directory", pathSplit[depth-1]))
			return data, "", code, "text/html", errors.New(fmt.Sprintf("no site data for: %v", pathSplit))
		}
		data, code := GetErrorSite(http.StatusNotFound, request.Host, request.URL.Path, "")
		return data, "", code, "text/html", errors.New(fmt.Sprintf("no site data for: %s", pathSplit))
	}
	data, encoding := file.data.getSmallest(availableEncodings)
	return data, encoding, 200, file.mimetype, nil
}

// CreateServe
//
// Registers a handle for '/' to serve the DefaultSite
func CreateServe() http.HandlerFunc {
	fun := func(w http.ResponseWriter, r *http.Request) {
		if settings.GetSettings().ServerOff.Get() {
			w.WriteHeader(http.StatusGone)
		}
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

		msg, encoding, code, mime, err := getSite(r, &availableEncodings)

		searchTime := time.Now()

		if err != nil {
			log.Err(err, fmt.Sprintf("Error getting site %s", r.URL.Path))
			w.WriteHeader(code)
		} else {
			w.Header().Set("Content-Type", mime)
			if encoding != "" {
				w.Header().Set("Content-encoding", string(encoding))
			}
		}

		_, er := w.Write(*msg)
		if er != nil {
			log.Err(er, "Error writing response:")
		}
		go LogAccess(code, int(time.Since(start).Microseconds()), int(searchTime.Sub(start).Microseconds()), err, er, r.Method, r.URL.Path, encoding)
	}

	return fun
}

type Encoding string

const (
	GZip    Encoding = "gzip"
	Deflate Encoding = "deflate"
	Brotli  Encoding = "br"
)
