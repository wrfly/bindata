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

	xs := fs
	for _, f := range fs {
		if !f.isDir {
			continue
		}
		// fill the dir
		for _, ff := range xs {
			if ff.dirP == f.p {
				f.dir = append(f.dir, ff.fileInfo)
				f.files = append(f.files, ff)
				f.assets = append(f.assets, ff)
			}
		}
	}

	return fs, nil
}

func Gen(dir string) (Assets, error) {
	if !filepath.IsAbs(dir) {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(pwd, dir)
	}

	fs, err := walk(dir, 0)
	if err != nil {
		return nil, err
	}

	d := &data{
		prefix: "/x",
		files:  make(map[string]*file, len(fs)),
	}
	all := &file{fileInfo: &fileInfo{isDir: true}}
	for _, f := range fs {
		name := filepath.Join(d.prefix, f.rlvP)
		if f.IsDir() && name != "/" {
			d.files[name+"/"] = f
		}
		d.files[name] = f

		all.files = append(all.files, f)
		all.dir = append(all.dir, f.fileInfo)
		all.assets = append(all.assets, f)
	}
	d.all = all

	return d, nil
}
