package main

import (
	"flag"
	"fmt"

	bindata "github.com/wrfly/bindata/lib"
)

var (
	packageName, resource, prefix, target string
)

func init() {
	flag.StringVar(&packageName, "package", "asset", "target package name")
	flag.StringVar(&resource, "resource", "", "resource dir")
	flag.StringVar(&prefix, "prefix", "/",
		"resource prefix, used for http server")
	flag.StringVar(&target, "target", "example/asset",
		"where to put the generated files")
	flag.Parse()
}

func main() {
	if resource == "" {
		flag.Usage()
		return
	}

	_, err := bindata.Gen(bindata.GenOption{
		Package:  packageName,
		Resource: resource,
		Prefix:   prefix,
		Target:   target,
	})
	if err != nil {
		fmt.Printf("err: %s", err)
	}
}
