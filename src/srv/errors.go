package srv

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type Errors uint16

func GetErrorSite(error Errors, host string, path string, additional string) (*[]byte, int) {
	var site string

	switch error {
	case http.StatusForbidden:
		site = "You are not allowed to access this URL."
	case http.StatusNotFound:
		site = "URL not found on server."
	case http.StatusMethodNotAllowed:
		site = "Method not allowed."
	case http.StatusInternalServerError:
		site = "An error happened while processing your Request."
	default:
		site = "Error not found"
	}

	site = fmt.Sprintf(`
<html>
	<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width"/>
		<title>%d | %s</title>
	</head>
	<body style="background:black;background:linear-gradient(140deg, rgb(7 10 15) 0%%, rgb(0,0,0) 50%%, rgb(7 9 10) 100%%);;color:white">
		<div style="margin:auto;width:50%%;padding:10px;position:absolute;bottom:50%%;right:50%%;transform:translate(50%%,50%%);overflow:hidden">
			<div style="display:flex;align-items:center;justify-content:space-between;gap:2em">
				<h1 style="margin-block:0.2em">%s</h1>
				<p>Error accessing %s</p>
			</div>		
			<div style="display:flex;align-items:center;justify-content:space-between;gap:2em">
				<p>%s</p>
				<p>%s</p>
				<button style="padding:8px 16px;color:white;border:white 1px solid;background:transparent;cursor:pointer;border-radius:1em" onclick="location.reload()">Reload</button>
			</div>
			<img src="https://http.cat/%d" style="width:80%%;margin-left: 10%%">
			<hr>
			<address>GoWebserver at %s running %s on %s</address>
		</div>
	</body>
</html>
`, error, http.StatusText(int(error)), http.StatusText(int(error)), path, site, additional, int(error), host, runtime.Version(), getOS())

	ret := []byte(site)

	return &ret, int(error)
}

func getOS() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "Not UNIX"
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "NAME=") {
			return strings.TrimSuffix(strings.TrimPrefix(s.Text(), "NAME=\""), "\"")
		}
	}
	return ""
}
