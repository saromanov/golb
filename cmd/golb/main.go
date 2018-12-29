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
	// ConfigСonsulKey defines key for load config from Consul
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
		ConfigСonsulKey = args[0]
	},
}

var configPathCmd = &cobra.Command{
	Use:   "config-path",
	Short: "Load config from path",
	Long:  "Load config from path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("key is not defined")
		}
		ConfigPath = args[0]
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(configConsulCmd)
	rootCmd.AddCommand(configPathCmd)

	var (
		cfg *config.Config
		err error
	)

	if ConfigPath != "" {
		cfg, err = config.ReadConfig(ConfigPath)
		if err != nil {
			panic(fmt.Sprintf("unable to read config: %v", err))
		}
	}
	if ConfigСonsulKey != "" {
		cfg, err = config.ReadConsulConfig(ConfigСonsulKey)
		if err != nil {
			panic(fmt.Sprintf("unable to read config: %v", err))
		}
	}
	if cfg == nil {
		cfg = config.MakeDefaultConfig()
	}

	g := golb.New(cfg)
	if err := g.Build(); err != nil {
		panic(err)
	}
	http.HandleFunc("/", g.HandleHTTP)
}
