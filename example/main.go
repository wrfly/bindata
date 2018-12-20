package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/wrfly/bindata"
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
	dir := "/home/mr/Documents/workspace/golang/src/github.com/wrfly/bindata/resource"
	asset, err := bindata.Gen(dir)
	if err != nil {
		panic(err)
	}
	http.Handle("/", asset)
	http.Handle("/x", asset)
	http.ListenAndServe(listen, nil)
}
