package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/saromanov/golb/config"
	"github.com/saromanov/golb/golb"
)

const defaultAddress = "127.0.0.1:8099"

func makeHTTPServer() {
	fmt.Println("Starting of the server...")
	err := http.ListenAndServe(defaultAddress, nil)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func makeHTTPSServer(cfg *config.Config) {
	err := http.ListenAndServeTLS(defaultAddress, cfg.CertFilePath, cfg.KeyFilePath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func main() {
	cfg, err := config.ReadConfig("./configs/config.json")
	if err != nil {
		panic(fmt.Sprintf("unable to read config: %v", err))
	}

	g := golb.New(cfg)
	g.Build()
	http.HandleFunc("/", g.HandleHTTP)
	if cfg.ServerScheme == "https" {
		makeHTTPSServer(cfg)
	}
	makeHTTPServer()
}
