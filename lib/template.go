package bindata

var fileTemplate = `
var %s = &file{
	fileInfo: &fileInfo{
		name:  "%s",
		isDir: %t,
		size:  %d,
		mode:  os.FileMode(%d),
		mTime: time.Unix(%d, 0),
		cType: "%s",
	},
	path:  "%s",
	dirP:  "%s",
	sPath: "%s",
	id:    %d,
	cb:    %s,
}
`

var fileBytesTemplate = `
var %s = []byte("%s")
`

var headerTemplate = `/*
CODE GENERATED BY "github.com/wrfly/bindata" 
@%s

Files:%s

DO NOT EDIT!
*/

package %s

import (
	"os"
	"time"
)
`

var footerTemplate = `
var fs = []*file{%s
}

var Data Assets

var d = &data{
	prefix: "%s",
	files:  make(map[string]*file, len(fs)),
}

func init() {
	for _, f := range fs {
		f.b = unCompress(f.cb)
		if !f.isDir || len(f.files) != 0 {
			continue
		}
		for _, ff := range fs {
			if ff.dirP == f.path {
				f.infos = append(f.infos, ff.fileInfo)
				f.files = append(f.files, ff)
				f.assets = append(f.assets, &fileReader{ff, 0})
			}
		}
	}

	all := &file{fileInfo: &fileInfo{isDir: true}}
	for _, f := range fs {
		if f.IsDir() {
			d.files[f.sPath+"/"] = f
		}
		d.files[f.sPath] = f
		d.files[f.path] = f
		all.files = append(all.files, f)
		all.infos = append(all.infos, f.fileInfo)
		all.assets = append(all.assets, &fileReader{f, 0})
	}
	d.all = all

	Data = d
}
`
