package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/saromanov/golb/config"
	"github.com/saromanov/golb/golb"
	"github.com/spf13/cobra"
)

var (
	// ConfigConsulKey defines key for load config from Consul
	ConfigСonsulKey string
	// ConfigPath defines path to load config
	ConfigPath string
)

var rootCmd = &cobra.Command{
	Use:   "golb",
	Short: "GoLB is a load balancer",
	Long:  `Implementation of the simple load balancer`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var configConsulCmd = &cobra.Command{
	Use:   "config-consul",
	Short: "Load config from Consul by the key",
	Long:  "Load config from Consul by the key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("key is not defined")
		}
		consulKey = args[0]
	},
}

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
	rootCmd.AddCommand(configConsulCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
