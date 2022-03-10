package main

import (
	"fmt"
	lg "log"
	"net/http"

	// graph "server/graphql"
	// gen "server/graphql/generated"

	"server/src/log"

	"server/src"
)

func main() {
	src.LoadConfig()
	log.Log("Loaded config:", fmt.Sprintf("%+v", src.GetConfig()))

	log.Log("Starting server")
	// log.DBInit()

	// err := serve.LoadSites()
	// if err != nil {
	// 	panic(err)
	// }

	serv := serve.CreateServe()

	webServer := &http.Server{Addr: ":" + fmt.Sprintf("%d", src.GetConfig().PortHTTPS), Handler: serv}
	webServer.ErrorLog = lg.New(&log.LogWriter{}, "", 0)
	go startWebServer(webServer, true)

	http.HandleFunc("/", graph.GetPlayground)
	http.Handle("/query", handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{Resolvers: graph.GenResolver()})))
	APIServer := &http.Server{Addr: ":" + fmt.Sprintf("%d", src.GetConfig().ApiPort), Handler: http.DefaultServeMux}
	APIServer.ErrorLog = lg.New(&log.LogWriter{}, "", 0)

	startAPI(APIServer)
}

func startWebServer(webServer *http.Server) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe Webserver with TLS started on localhost%s", webServer.Addr))
	err := webServer.ListenAndServeTLS(src.CertsFile, src.KeyFile)

	if err != nil {
		log.Err(err, "Error starting webServer")
	}
}

func startAPI(api *http.Server) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe API with TLS started on localhost%s", api.Addr))
	err := api.ListenAndServeTLS(src.CertsFile, src.KeyFile)

	if err != nil {
		log.Err(err, "Error starting Api")
	}
}
