package srv

import (
	"errors"
	"fmt"
	"html"
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
	if uint16(len(request.URL.String())) > settings.GetSettings().MaxURILength.Get() {
		data, code := GetErrorSite(http.StatusRequestURITooLong, request.Host, request.URL.Path, "")
		return data, "", code, "text/html", errors.New(fmt.Sprintf("URI to long (%v)", len(request.URL.String())))
	}

	url := html.EscapeString(request.URL.Path)
	host := html.EscapeString(request.Host)

	for _, forbidden := range settings.GetSettings().Forbidden.Get() {
		switch forbidden.Type {
		case settings.FileExtension:
			if strings.HasSuffix(url, "."+forbidden.Data) {
				data, code := GetErrorSite(http.StatusForbidden, host, url, "forbidden by FileExtension")
				return data, "", code, "text/html", errors.New(fmt.Sprintf("forbidden by extension: %s (%s)", forbidden.Data, url))
			}
		case settings.AbsoluteFile:
			if url == forbidden.Data {
				data, code := GetErrorSite(http.StatusForbidden, host, url, "forbidden by absolute path")
				return data, "", code, "text/html", errors.New(fmt.Sprintf("forbidden by absolute path: %s", forbidden.Data))
			}
		case settings.AbsoluteDirectory:
			if strings.HasPrefix(url, forbidden.Data) {
				data, code := GetErrorSite(http.StatusForbidden, host, url, "forbidden by absolute DirPath")
				return data, "", code, "text/html", errors.New(fmt.Sprintf("forbidden by absolute DirPath: %s (%s)", forbidden.Data, url))
			}
		case settings.Regex:
			if forbidden.Regex.Match([]byte(url)) {
				data, code := GetErrorSite(http.StatusForbidden, host, url, "forbidden by regex")
				return data, "", code, "text/html", errors.New(fmt.Sprintf("forbidden by regex: %s (%s)", forbidden.Data, url))
			}
		}
	}

	pathSplit := strings.Split(url, "/")[1:]

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
			data, code := GetErrorSite(http.StatusNotFound, host, url, fmt.Sprintf("%s is no file, but a directory", pathSplit[depth-1]))
			return data, "", code, "text/html", errors.New(fmt.Sprintf("no site data for: %v", pathSplit))
		}
		data, code := GetErrorSite(http.StatusNotFound, host, url, "")
		return data, "", code, "text/html", errors.New(fmt.Sprintf("no site data for: %s", pathSplit))
	}
	data, encoding := file.data.getSmallest(availableEncodings)
	return data, encoding, 200, file.mimetype, nil
}

func CreateServe(loading *bool) http.HandlerFunc {
	fun := func(w http.ResponseWriter, r *http.Request) {
		if *loading {
			data := GetLoadingSite(r.Host, r.URL.Path)

			w.WriteHeader(http.StatusTeapot)
			w.Header().Set("Content-Type", "text/html")
			_, er := w.Write(*data)
			if er != nil {
				log.Err(er, "Error writing response:")
			}

			return
		}

		if settings.GetSettings().ServerOff.Get() {
			w.WriteHeader(http.StatusGone)
			return
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

		// get actual site
		searchTime := time.Now()
		msg, encoding, code, mime, err := getSite(r, &availableEncodings)

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
