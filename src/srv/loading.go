package srv

import (
	"fmt"
	"runtime"
)

func GetLoadingSite(host string, path string) *[]byte {
	site := fmt.Sprintf(loadingSite, path, host, runtime.Version(), getOS())

	ret := []byte(site)

	return &ret
}

const loadingSite = `
<html>
	<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width"/>
		<title>Server Loading</title>
	</head>
	<body style="background:black;background:linear-gradient(140deg, rgb(7 10 15) 0%%, rgb(0,0,0) 50%%, rgb(7 9 10) 100%%);color:white">
		<div style="margin:auto;width:50%%;padding:10px;position:absolute;bottom:50%%;right:50%%;transform:translate(50%%,50%%);overflow:hidden;display:flex;flex-direction:column">
			<div style="display:flex;align-items:center;justify-content:space-between;gap:2em">
				<h1 style="margin-block:0.2em;flex-shrink:0">Server is still starting</h1>
				<p>Error accessing %s</p>
			</div>		
			<div style="display:flex;align-items:center;justify-content:space-between;gap:2em">
				<p>-- Progress --</p>
				<button style="padding:8px 16px;color:white;border:white 1px solid;background:transparent;cursor:pointer;border-radius:1em" onclick="location.reload()">Reload</button>
			</div>
			<img src="https://media.tenor.com/RVvnVPK-6dcAAAAC/reload-cat.gif" style="align-self:center;max-height:80vh">
			<hr style="width:100%%">
			<address>GoWebserver at %s running %s on %s</address>
		</div>
	</body>
</html>
`
