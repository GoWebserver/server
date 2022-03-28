package main

import (
	"fmt"
	lg "log"
	"net/http"

	// graph "server/graphql"
	// gen "server/graphql/generated"

	"server/src/config"
	"server/src/log"
	"server/src/srv"

	"server/src"
)

func main() {
	config.LoadConfig()
	log.Log("Loaded config:", fmt.Sprintf("%+v", config.GetConfig()))

	config.LoadSettings()

	log.Log("Starting server")
	src.DBInit()

	srv.LoadMime()
	err := srv.LoadSites()
	if err != nil {
		panic(err)
	}

	serv := srv.CreateServe()

	webServer := &http.Server{Addr: ":" + fmt.Sprintf("%d", config.GetConfig().Port), Handler: serv}
	webServer.ErrorLog = lg.New(&log.LogWriter{}, "", 0)
	startWebServer(webServer)

	// http.HandleFunc("/", graph.GetPlayground)
	// http.Handle("/query", handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{Resolvers: graph.GenResolver()})))
	// APIServer := &http.Server{Addr: ":" + fmt.Sprintf("%d", config.GetConfig().ApiPort), Handler: http.DefaultServeMux}
	// APIServer.ErrorLog = lg.New(&log.LogWriter{}, "", 0)
	//
	// startAPI(APIServer)
}

func startWebServer(webServer *http.Server) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe Webserver with TLS started on localhost%s", webServer.Addr))
	err := webServer.ListenAndServeTLS(config.CertsFile, config.KeyFile)

	if err != nil {
		log.Err(err, "Error starting webServer")
	}
}

func startAPI(api *http.Server) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe API with TLS started on localhost%s", api.Addr))
	err := api.ListenAndServeTLS(config.CertsFile, config.KeyFile)

	if err != nil {
		log.Err(err, "Error starting Api")
	}
}
