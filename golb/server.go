// Creating of server for sending requests and getting stats
// from load balancer
package golb

import (
	"fmt"
	"log"
	"net/http"

	"github.com/saromanov/golb/config"
)

const defaultAddress = "127.0.0.1:8099"

// Server provides definition for stat server
type Server struct {
}

func makeHTTPServer() {
	fmt.Println("Starting of the server...")
	err := http.ListenAndServe(defaultAddress, nil)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func makeHTTPSServer(cfg *config.Config) {
	fmt.Println("Starting of the server")
	err := http.ListenAndServeTLS(defaultAddress, cfg.CertFilePath, cfg.KeyFilePath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
