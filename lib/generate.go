package bindata

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type GenOption struct {
	Package  string // package name
	Prefix   string // file prefix
	Target   string // where to put the generated file
	Resource string // a single file or a dir

	AssetName   string // default=asset.go
	BindataName string // default=bindata.go
}

func walk(root string) (fs []*file, err error) {
	id := 0
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("get file [%s] error, info is nil", path)
		}
		xPath := filepath.Join("/", strings.TrimPrefix(path, root))
		got := &file{
			fileInfo: &fileInfo{
				name:  info.Name(),
				isDir: info.IsDir(),
				size:  info.Size(),
				mode:  info.Mode(),
				mTime: info.ModTime(),
				cType: mime.TypeByExtension(filepath.Ext(path)),
			},
			path: xPath,
			dirP: filepath.Dir(xPath),
			id:   id,
		}
		id++

		if xPath == "/" {
			got.name = "/"
		}

		if !info.IsDir() {
			got.b, err = ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read file %s error: %s", path, err)
			}
		}
		fs = append(fs, got)

		return nil
	})
	if err != nil {
		return
	}
	fill(fs)

	return
}

func fill(fs []*file) {
	xs := fs
	for _, f := range fs {
		if !f.isDir {
			continue
		}
		// fill the dir
		for _, ff := range xs {
			if ff.dirP == f.path {
				f.infos = append(f.infos, ff.fileInfo)
				f.files = append(f.files, ff)
				f.assets = append(f.assets, &fileReader{ff, 0})
			}
		}
	}
}

func Gen(opts GenOption) (*data, error) {
	// validate options
	if opts.Package == "" {
		opts.Package = "asset"
	}
	if opts.Prefix == "" {
		opts.Prefix = "/"
	}
	if opts.Target == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		opts.Target = filepath.Join(wd, "asset")
	}

	// make data
	d, err := buildData(opts.Resource, opts.Prefix)
	if err != nil {
		return nil, err
	}

	// make writerTo
	w, err := buildWriter(d, opts.Prefix, opts.Package)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(opts.Target)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(opts.Target, 0755); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else if !info.IsDir() {
		return nil, fmt.Errorf("target is not a directory")
	}

	if opts.AssetName == "" {
		opts.AssetName = "asset.go"
	}

	targetPkg := filepath.Join(opts.Target, opts.AssetName)
	f, err := os.Create(targetPkg)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = w.WriteTo(f)
	if err != nil {
		return nil, err
	}

	if opts.BindataName == "" {
		opts.BindataName = "bindata.go"
	}

	targetBin := filepath.Join(opts.Target, opts.BindataName)
	fBindata, err := os.Create(targetBin)
	if err != nil {
		return nil, err
	}
	defer fBindata.Close()
	_, err = fBindata.WriteString(fmt.Sprintf(bindataTemplate, opts.Package))

	return d, err
}

func buildData(resource, prefix string) (*data, error) {
	if !filepath.IsAbs(resource) {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		resource = filepath.Join(pwd, resource)
	}

	fs, err := walk(resource)
	if err != nil {
		return nil, err
	}

	var size int64
	all := &file{fileInfo: &fileInfo{isDir: true}}
	for _, f := range fs {
		size += f.size
		if f.path == "/" {
			f.sPath = filepath.Join(prefix, f.dirP)
		} else {
			f.sPath = filepath.Join(prefix, f.dirP, f.name)
		}
		all.files = append(all.files, f)
		all.infos = append(all.infos, f.fileInfo)
		all.assets = append(all.assets, &fileReader{f, 0})
	}

	return &data{
		prefix: prefix,
		files:  make(map[string]*file, len(fs)),
		all:    all,
	}, nil
}

func buildWriter(d *data, prefix, pName string) (io.WriterTo, error) {
	w := new(bytes.Buffer)

	// package header
	filesStr := ""
	for _, f := range d.all.files {
		filesStr = fmt.Sprintf("%s\n\t%s", filesStr, f.path)
	}
	w.WriteString(fmt.Sprintf(headerTemplate, time.Now().Format(time.RFC3339), filesStr, pName))

	// files
	var (
		wg sync.WaitGroup
		m  sync.Mutex
	)
	wg.Add(len(d.all.files))
	for _, f := range d.all.files {
		go func(f *file) {
			content := printFile(f.sPath, f)
			m.Lock()
			w.WriteString(content)
			m.Unlock()
			wg.Done()
		}(f)
	}
	wg.Wait()

	// package footer
	names := []string{}
	for _, f := range d.all.files {
		names = append(names, f.keyFileName())
	}
	fs := ""
	for i, n := range names {
		if i%5 == 0 {
			fs += "\n\t\t"
		} else {
			fs += " "
		}
		fs += n + ","
	}
	w.WriteString(fmt.Sprintf(footerTemplate, fs, prefix))

	return w, nil
}

func printFile(name string, f *file) string {
	compressedBytes := compress(f.b)
	// print bytes
	bsW := bytes.NewBuffer(make([]byte, 0, 1e3))
	for i, b := range compressedBytes {
		if i%15 == 0 && len(compressedBytes) > 15 {
			bsW.WriteString(fmt.Sprintf("\" +\n\t\""))
		}
		bsW.WriteString(fmt.Sprintf("\\x%02x", b))
	}

	contentW := bytes.NewBuffer(make([]byte, 0, bsW.Len()))

	contentW.WriteString(fmt.Sprintf(fileBytesTemplate, f.keyBytesName(), bsW))

	// print file
	contentW.WriteString(fmt.Sprintf(fileTemplate,
		f.keyFileName(),
		f.name,
		f.isDir,
		f.size,
		f.mode,
		f.mTime.Unix(),
		f.cType,
		f.path,
		f.dirP,
		f.sPath,
		f.id,
		f.keyBytesName(),
	))

	return contentW.String()
}

func compress(in []byte) []byte {
	zipBuffer := bytes.NewBuffer(make([]byte, 0, 1e3))
	zipWriter := zlib.NewWriter(zipBuffer)
	zipWriter.Write(in)
	zipWriter.Close()
	return zipBuffer.Bytes()
}
