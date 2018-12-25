package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	embeddedData := asset.Data

	// list files
	assets := embeddedData.List()
	for _, file := range assets {
		info, _ := file.File()
		stat, _ := info.Stat()
		fmt.Printf("contains file: %s [%d]\n",
			stat.Name(), stat.Size())
	}

	// open one file
	n := "/sub/sub_main.css"
	f, _ := embeddedData.Open(n)
	bs, _ := ioutil.ReadAll(f)
	fmt.Printf("file %s's data: %s\n", n, bs)

	// serve http
	http.Handle("/", embeddedData)
	http.Handle("/html", embeddedData)
	http.ListenAndServe(listen, nil)
}
