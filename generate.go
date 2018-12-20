package bindata

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		if f.isDir {
			for _, ff := range fs {
				if !ff.isDir && ff.dirP == f.p {
					f.files = append(f.files, ff)
				}
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
	for _, f := range fs {
		if f.isDir {
			fmt.Println("->", f.p)
			for _, ff := range f.files {
				fmt.Println(f.p, ":", ff.name)
			}
		}
	}

	d := &data{
		prefix: "/",
		files:  make(map[string]*file, 0),
	}
	d.files["/"] = &file{
		fileInfo: &fileInfo{
			name: "/",
			size: 4,
		},
		b: []byte("1234"),
	}
	d.files["/x"] = &file{
		fileInfo: &fileInfo{
			name: "",
			size: 3,
		},
		b: []byte("123"),
	}
	return d, nil
}
