package bindata

import "testing"

func TestGenerate(t *testing.T) {
	path := "resource"
	fs, err := Gen(path)
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
