package lib

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	fs, err := Gen(Option{
		Package:  "asset",
		Resource: "../resource",
		Prefix:   "/html",
		Target:   "../example/asset",
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fs.List() {
		stat, err := f.Stat()
		if err != nil {
			t.Error(err)
			continue
		}
		t.Logf("got file: %s", stat.Name())
	}

}
