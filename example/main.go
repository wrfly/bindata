package main

import (
	"bytes"
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
	// list files
	assets := asset.List()
	for _, file := range assets {
		stat, _ := file.Stat()
		fmt.Printf("contains file: %s [%d]\n",
			stat.Name(), stat.Size())
	}

	// open one file
	n := "/sub/sub_main.css"
	f, err := asset.Find(n)
	if err != nil {
		panic(err)
	}
	bs, _ := ioutil.ReadAll(f)
	fmt.Printf("file %s's data: %s\n", n, bs)

	t := f.Template()
	w := bytes.NewBuffer(nil)
	t.Execute(w, map[string]interface{}{
		"hey": "girl",
	})
	fmt.Printf("template: %s", w)

	// serve http
	http.HandleFunc("/", asset.Handler)
	http.HandleFunc("/html", asset.Handler)
	http.ListenAndServe(listen, nil)
}
