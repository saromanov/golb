package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/saromanov/golb/config"
	"github.com/saromanov/golb/golb"
	"github.com/saromanov/golb/server"
)

const DefaultAddress = "127.0.0.1:8099"

// makeDefault provides creating of the
// Golb object with default params
func makeDefault() golb.GoLB {
	return golb.GoLB{
		MaxConnections:    10,
		ClientTimeout:     1 * time.Second,
		ConnectionTimeout: 20 * time.Second,
		Protocol:          "tcp",
		Scheme:            "http",
		Servers: []*server.Server{
			&server.Server{
				Host: "127.0.0.1",
				Port: 8900,
			},
			&server.Server{
				Host: "127.0.0.1",
				Port: 8901,
			},
		},
	}
}
func main() {

	g := makeDefault()
	cfg, err := config.ReadConfig("./configs/config.json")
	if err == nil {
		g = config.MakeGoLBObject(cfg)
	}
	g.Build()

	http.HandleFunc("/", g.HandleHTTP)
	fmt.Println("Starting of the server...")
	err = http.ListenAndServe(DefaultAddress, nil)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}
