package bindata

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const maxWalkDepth = 10

var (
	errMaxDepthExceeded = errors.New("max depth is 10")
)

var pf = fmt.Printf

func walk(root string, depth int) ([]*file, error) {
	if depth > maxWalkDepth {
		return nil, errMaxDepthExceeded
	}
	fs := []*file{}
	if err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			got := &file{
				fileInfo: &fileInfo{
					name:  info.Name(),
					isDir: info.IsDir(),
					size:  info.Size(),
					mode:  info.Mode(),
					mTime: info.ModTime(),
				},
				p:    path,
				dirP: filepath.Dir(path),
				rlvP: strings.TrimPrefix(path, root),
			}

			if !info.IsDir() {
				got.b, err = ioutil.ReadFile(path)
				if err != nil {
					return err
				}
			}
			fs = append(fs, got)
			return nil
		}); err != nil {
		return fs, err
	}

	for _, f := range fs {
		if !f.isDir {
			continue
		}
		// fill the dir
		for _, ff := range fs {
			if !ff.isDir && ff.dirP == f.p {
				f.dir = append(f.dir, ff.fileInfo)
				f.files = append(f.files, ff)
				f.assets = append(f.assets, ff)
			}
		}
	}

	return fs, nil
}

func Gen(dir string) (Assets, error) {
	fs, err := walk(dir, 0)
	if err != nil {
		return nil, err
	}

	d := &data{
		prefix: "/x",
		files:  make(map[string]*file, len(fs)),
	}
	for _, f := range fs {
		name := filepath.Join(d.prefix, f.rlvP)
		if f.IsDir() {
			if name != "/" {
				d.files[name+"/"] = f
			}
		}
		d.files[name] = f
	}
	return d, nil
}
