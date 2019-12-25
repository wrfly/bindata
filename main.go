package main

import (
	l "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wrfly/gua"

	bindata "github.com/wrfly/bindata/lib"
)

type config struct {
	PkgName string `name:"pkg" desc:"package name" default:"github.com/wrfly/bindata/example/asset"`
	Source  string `name:"src" desc:"resource dir"`
	Dest    string `name:"dest" desc:"path to store go asset files, default is the pkg path"`
	Prefix  string `name:"prefix" desc:"resource prefix, used for static files in HTTP server" default:"/"`

	WithTime bool `name:"with-time" desc:"generate file with its time info"`
	WithMod  bool `name:"with-mod" desc:"generate file with its mod info"`
}

var log = l.New(os.Stdout, "", 0)

func main() {
	cfg := new(config)
	if err := gua.Parse(cfg); err != nil {
		log.Printf("err: %s", err)
		return
	}

	var pkg string
	if strings.Contains(cfg.PkgName, "/") {
		x := strings.Split(cfg.PkgName, "/")
		pkg = x[len(x)-1]
	}

	// set default cfg.Dest path as package path
	if cfg.Dest == "" {
		cfg.Dest = filepath.Join(os.Getenv("GOPATH"), "src", cfg.PkgName)
	}

	log.Printf("package: [%s]", cfg.PkgName)
	log.Printf("prefix:  [%s]", cfg.Prefix)
	log.Printf("src:     [%s]", cfg.Source)
	log.Printf("dest:    [%s]", cfg.Dest)

	start := time.Now()
	_, err := bindata.Gen(bindata.GenOption{
		Package:  pkg,
		Resource: cfg.Source,
		Prefix:   cfg.Prefix,
		Target:   cfg.Dest,
	})
	if err != nil {
		log.Printf("err: %s", err)
		return
	}

	info, err := os.Stat(filepath.Join(cfg.Dest, "asset.go"))
	if err != nil {
		log.Printf("err: %s", err)
		return
	}

	log.Printf("done; size=%dK, use=%s",
		info.Size()/1024, time.Since(start))

}
