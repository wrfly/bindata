package main

import (
	"github.com/wrfly/bindata"
)

func main() {
	_, err := bindata.Gen(bindata.GenOption{
		Package:  "main",
		Resource: "../resource",
		Prefix:   "/html",
		Target:   "../example",
	})
	if err != nil {
		panic(err)
	}

}
