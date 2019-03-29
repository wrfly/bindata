package main

import (
	"flag"
	"fmt"
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

	fmt.Printf("package=[%s]\nprefix=[%s]\ntarget=[%s]\nresource=[%s]\n",
		packageName, prefix, target, resource)
	fmt.Println("... ...")

	start := time.Now()
	_, err := bindata.Gen(bindata.GenOption{
		Package:  pkg,
		Resource: resource,
		Prefix:   prefix,
		Target:   target,
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	info, _ := os.Stat(filepath.Join(target, "asset.go"))

	fmt.Printf("done, size=%dK use=%s\n",
		info.Size()/1024, time.Since(start))

}
