// Creating of server for sending requests and getting stats
// from load balancer
package golb

import (
	"fmt"
	"log"
	"net/http"

	"github.com/saromanov/golb/config"
)

const defaultMetricsServer = "127.0.0.1:8099"

// Server provides definition for stat server
type Server struct {
}

func makeHTTPMetricsServer() {
	fmt.Println("Starting of the metrics server...")
	err := http.ListenAndServe(defaultMetricsServer, nil)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func makeHTTPSMetricsServer(cfg *config.Config) {
	fmt.Println("Starting of the server")
	err := http.ListenAndServeTLS(defaultMetricsServer, cfg.CertFilePath, cfg.KeyFilePath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
