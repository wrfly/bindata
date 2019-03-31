package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	bindata "github.com/wrfly/bindata/lib"
)

var (
	packageName, resource, prefix, target string
)

func init() {
	flag.StringVar(&resource, "resource", "", "resource dir")
	flag.StringVar(&packageName, "pkg",
		"github.com/wrfly/bindata/example/asset", "target package name")
	flag.StringVar(&prefix, "prefix",
		"/", "resource prefix, used for http server")
	flag.StringVar(&target, "target",
		"", "where to put the generated files, default is the package's path")

	flag.Parse()
}

func main() {
	if resource == "" {
		flag.Usage()
		return
	}

	var pkg string
	if strings.Contains(packageName, "/") {
		x := strings.Split(packageName, "/")
		pkg = x[len(x)-1]
	}

	// set default target path as package path
	if target == "" {
		gopath := os.Getenv("GOPATH")
		target = filepath.Join(gopath, "src", packageName)
	}

	log.Printf("package: [%s]", packageName)
	log.Printf("prefix: [%s]", prefix)
	log.Printf("target: [%s]", target)
	log.Printf("resource: [%s]", resource)

	start := time.Now()
	_, err := bindata.Gen(bindata.GenOption{
		Package:  pkg,
		Resource: resource,
		Prefix:   prefix,
		Target:   target,
	})
	if err != nil {
		log.Printf("err: %s\n", err)
		return
	}

	info, err := os.Stat(filepath.Join(target, "asset.go"))
	if err != nil {
		log.Printf("err: %s\n", err)
		return
	}

	log.Printf("done; size=%dK, use=%s\n",
		info.Size()/1024, time.Since(start))

}
