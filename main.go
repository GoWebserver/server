package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/lucas-clemente/quic-go/http3"

	"server/src/config"
	"server/src/log"
	"server/src/settings"
	"server/src/srv"

	"server/src"
)

var loading = true

func main() {
	config.LoadConfig()
	log.Log("Loaded config:", fmt.Sprintf("%+v", config.GetConfig()))

	settings.LoadDefaultSettings()

	log.Log("Starting server")
	src.DBInit()

	go func() {
		loading = true
		srv.LoadSites()
		loading = false
	}()

	serv := srv.CreateServe(&loading)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		startWebServer("localhost:"+fmt.Sprintf("%d", config.GetConfig().Port), serv)
		wg.Done()
	}()

	// http.HandleFunc("/", graph.GetPlayground)
	// http.Handle("/query", handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{Resolvers: graph.GenResolver()})))
	// APIServer := &http.Server{Addr: ":" + fmt.Sprintf("%d", config.GetConfig().ApiPort), Handler: http.DefaultServeMux}
	// APIServer.ErrorLog = lg.New(&log.LogWriter{}, "", 0)
	//
	// startAPI(APIServer)

	wg.Wait()
}

func startWebServer(addr string, handler http.Handler) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe Webserver HTTP/3 with TLS started on https://%s", addr))
	err := http3.ListenAndServe(addr, config.CertsFile, config.KeyFile, handler)

	if err != nil {
		log.Err(err, "Error starting webServer")
		panic(err)
	}
}

/*
func startAPI(api *http.Server) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe API with TLS started on localhost%s", api.Addr))
	err := api.ListenAndServeTLS(config.CertsFile, config.KeyFile)

	if err != nil {
		log.Err(err, "Error starting Api")
		panic(err)
	}
}
*/
