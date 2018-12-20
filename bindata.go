package bindata

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Assets interface {
	Asset(name string) (Asset, error)
	Open(name string) (http.File, error)              // implement http.FileSystem
	ServeHTTP(w http.ResponseWriter, r *http.Request) // implement http.FileServer
}

type Asset interface {
	List() ([]Asset, error)
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

var (
	errSeekInvalid  = errors.New("invalid whence")
	errSeekNegative = errors.New("negative position")
	errNotDir       = errors.New("file is not a dir")
)

type data struct {
	prefix string
	files  map[string]*file
}

func (d *data) Open(name string) (http.File, error) {
	if f, found := d.files[name]; found {
		return &fileReader{f, 0}, nil
	}
	return nil, os.ErrNotExist
}

func (d *data) Asset(name string) (Asset, error) {
	if f, found := d.files[name]; found {
		return f, nil
	}
	return nil, os.ErrNotExist
}

func (d *data) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f, found := d.files[r.RequestURI]; found {
		w.Write(f.b)
		w.Header().Set("Content-Length", fmt.Sprint(f.size))
		w.Header().Set("Content-Type", f.cType)
		w.Header().Set("Date", fmt.Sprint(f.mTime))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	return
}

type fileReader struct {
	*file
	i int64
}

func (r *fileReader) Read(p []byte) (n int, err error) {
	if r.i >= r.size {
		return 0, io.EOF
	}
	n = copy(p, r.file.b[r.i:])
	r.i += int64(n)
	return
}

func (r *fileReader) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = r.size + offset
	default:
		return 0, errSeekInvalid
	}
	if abs < 0 {
		return 0, errSeekNegative
	}
	r.i = abs
	return abs, nil
}

func (f *fileReader) Close() error {
	return nil
}

type file struct {
	*fileInfo
	p      string // path
	b      []byte // data
	dir    []os.FileInfo
	files  []*file
	assets []Asset
	dirP   string // dir path
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, errors.New("not dir")
	}
	if count < 0 {
		return nil, nil
	}
	if count >= len(f.dir) {
		count = len(f.dir) - 1
	}
	return f.dir[:count], nil
}

func (f *file) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *file) List() ([]Asset, error) {
	if f.isDir {
		return f.assets, nil
	}
	return nil, errNotDir
}

type fileInfo struct {
	name  string
	isDir bool
	size  int64
	mode  os.FileMode
	mTime time.Time
	cType string
	// cTime time.Time
}

// base name of the file
func (f *fileInfo) Name() string {
	return f.name
}

// length in bytes for regular files; system-dependent for others
func (f *fileInfo) Size() int64 {
	return f.size
}

// file mode bits
func (f *fileInfo) Mode() os.FileMode {
	return f.mode
}

// modification time
func (f *fileInfo) ModTime() time.Time {
	return f.mTime
}

// abbreviation for Mode().IsDir()
func (f *fileInfo) IsDir() bool {
	return f.isDir
}

// underlying data source (can return nil)
func (f *fileInfo) Sys() interface{} {
	return nil
}
