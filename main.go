package main

import (
	"crypto/tls"
	"fmt"
	lg "log"
	"net/http"
	"sync"

	// graph "server/graphql"
	// gen "server/graphql/generated"

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

	wg := sync.WaitGroup{}
	wg.Add(1)

	crt, err := tls.LoadX509KeyPair(config.CertsFile, config.KeyFile)
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{crt},
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	webServer := &http.Server{
		Addr:      ":" + fmt.Sprintf("%d", config.GetConfig().Port),
		Handler:   serv,
		TLSConfig: tlsConfig,
		ErrorLog:  lg.New(&log.LogWriter{}, "", 0),
	}

	go func() {
		startWebServer(webServer)
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

func startWebServer(webServer *http.Server) {
	// blocks if success
	log.Log(fmt.Sprintf("ListenAndServe Webserver with TLS started on https://localhost%s", webServer.Addr))
	err := webServer.ListenAndServeTLS("", "") // files get ignored, already provided via tlsConfig

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
