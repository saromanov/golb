// Creating of server for sending requests and getting stats
// from load balancer
package golb

// Server provides definition for stat server
type Server struct {

}


func createTLSServer(crt, key string) error {
    err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}