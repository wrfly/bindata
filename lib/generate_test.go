package bindata

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	fs, err := Gen(GenOption{
		Package:  "main",
		Resource: "resource",
		Prefix:   "/html",
		Target:   "example",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fs.List() {
		info, err := f.File()
		if err != nil {
			t.Error(err)
			continue
		}
		stat, err := info.Stat()
		if err != nil {
			t.Error(err)
			continue
		}
		t.Logf("got file: %s", stat.Name())
	}

}
