package main

import (
	"flag"
	"fmt"
	"net/http"
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
	http.Handle("/", Data)
	http.Handle("/html", Data)
	http.ListenAndServe(listen, nil)
}
