package main

import (
	"fmt"
	"sync"

	"github.com/lucas-clemente/quic-go/http3"

	"server/src/config"
	"server/src/log"
	"server/src/settings"
	"server/src/srv"

	"server/src"
)

func main() {
	config.LoadConfig()
	log.Log("Loaded config:", fmt.Sprintf("%+v", config.GetConfig()))

	settings.LoadDefaultSettings()

	log.Log("Starting server")
	src.DBInit()

	srv.LoadSites()

	serv := srv.CreateServe()

	webServer3 := &http3.Server{
		Addr:    "localhost:" + fmt.Sprintf("%d", config.GetConfig().Port),
		Handler: serv,
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		startWebServer(webServer3)
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

func startWebServer(webServer *http3.Server) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe Webserver HTTP/3 with TLS started on https://%s", webServer.Addr))
	// err := webServer.ListenAndServeTLS(config.CertsFile, config.KeyFile)
	err := http3.ListenAndServe(webServer.Addr, config.CertsFile, config.KeyFile, webServer.Handler)

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
