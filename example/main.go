package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/wrfly/bindata/example/asset"
)

var (
	port   = flag.Int("port", 8080, "listen port")
	listen string
)

func init() {
	flag.Parse()
	listen = fmt.Sprintf(":%d", *port)
}

func main() {
	http.Handle("/", asset.Data)
	http.Handle("/html", asset.Data)
	http.ListenAndServe(listen, nil)
}
