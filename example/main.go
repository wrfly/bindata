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

	for _, x := range asset.List() {
		f, err := x.File()
		if err != nil {
			panic(err)
		}
		info, err := f.Stat()
		if err != nil {
			panic(err)
		}
		fmt.Println(info.Name(), info.IsDir())
	}

	http.Handle("/", asset)
	http.Handle("/x", asset)
	http.ListenAndServe(listen, nil)
}
