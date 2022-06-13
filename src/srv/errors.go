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
		site = getForbidden(additional)
	case http.StatusNotFound:
		site = getNotFound(additional)
	case http.StatusMethodNotAllowed:
		site = getMethodNotAllowed(additional)
	case http.StatusInternalServerError:
		site = getInternalServerError(additional)
	default:
		site = getErrorNotFound(additional)
	}

	site = fmt.Sprintf(`
<html>
	<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width"/>
		<title>%d | %s</title>
	</head>
	<body style="background:black;color:white">
		<div style="margin:auto;width:50%%;padding:10px;position:absolute;bottom:50%%;right:50%%;transform:translate(50%%,50%%)">
			<div style="display:flex;align-items:center;justify-content:space-between;">
				<h1 style="margin-block:0.2em">%s</h1>
				<p>Error accessing %s</p>
			</div>		
			<div style="display:flex;align-items:center;justify-content:space-between">
				%s
				<button style="padding:8px 16px;color:white;border:white 1px solid;background:transparent;cursor:pointer" onclick="location.reload()">Reload</button>
			</div>
			<img src="https://http.cat/%d" style="width:100%%">
			<hr>
			<address>GoWebServer at %s running %s on %s</address>
		</div>
	</body>
</html>
`, error, http.StatusText(int(error)), http.StatusText(int(error)), path, site, int(error), host, runtime.Version(), getOS())

	ret := []byte(site)

	return &ret, int(error)
}

func getForbidden(additional string) string {
	return fmt.Sprintf(`
<p>You are not allowed to access this site.</p>
<p>%s</p>
	`, additional)
}

func getNotFound(additional string) string {
	return fmt.Sprintf(`
<p>Site not found on server.</p>
<p>%s</p>
	`, additional)
}

func getMethodNotAllowed(additional string) string {
	return fmt.Sprintf(`
<p>Method not allowed.</p>
<p>%s</p>
	`, additional)
}

func getInternalServerError(additional string) string {
	return fmt.Sprintf(`
<p>An error happend while processing your Request</p>
<p>%s</p>
	`, additional)
}

func getErrorNotFound(additional string) string {
	return fmt.Sprintf(`
<p>Error not found</p>
<p>%s</p>
	`, additional)
}

func getOS() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return ""
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "NAME") {
			return strings.TrimPrefix(s.Text(), "NAME=")
		}
	}
	return ""
}
