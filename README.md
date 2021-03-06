# BinData

Although there are many golang programs who transfer files into go
packages, such as [jteeuwen/go-bindata](https://github.com/jteeuwen/go-bindata),
[a-urth/go-bindata](https://github.com/a-urth/go-bindata), and
[go-bindata-assetfs](https://github.com/elazarl/go-bindata-assetfs), but I still
felt not good with them.

So I wrote one myself.

[![GoDoc](https://godoc.org/github.com/wrfly/bindata/lib?status.svg)](https://godoc.org/github.com/wrfly/bindata/lib)

## Install

Install the `bindata` cmd with:

```bash
go get github.com/wrfly/bindata
```

## Usage

```txt
Usage of ./bindata:
 -dest        path to store go asset files, default is the pkg path
 -pkg         package name [github.com/wrfly/bindata/example/asset]
 -prefix      resource prefix, used for static files in HTTP server [/]
 -src         resource dir
 -with-mod    generate file with its mod info [false]
 -with-time   generate file with its time info [false]
```

### Generate embedded data

For example, if you want to generate a asset package, who's path
is `github.com/wrfly/bindata/example/asset`, the resource files are
located in `resource/` directory, then yoy can use the command below:

```bash
bindata -pkg github.com/wrfly/bindata/example/asset \
    -src "resource/"
```

After execute the command:

```txt
$ ls $GOPATH/src/github.com/wrfly/bindata/example/asset
asset.go  bindata.go
```

And now, you can use this `asset` package with:

```golang
package main

import (
    "net/http"

    "github.com/wrfly/bindata/example/asset"
)

func main() {
    http.HandleFunc("/", asset.Handler)
    http.ListenAndServe(":8080", nil)
}
```

Look at the [example](example) folder for more details and usage examples.