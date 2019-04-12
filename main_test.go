package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/dchest/safefile"
	"github.com/m-lab/gcs-downloader/atomicfile"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func Test_main(t *testing.T) {
	dir, err := ioutil.TempDir("", "gcs-downloader-main-")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		source     string
		atomicFile atomicfile.File
	}{
		{
			name:       "ok",
			source:     "gs://m-lab-gcs-downloader-mlab-testing/t1/okay.txt",
			atomicFile: atomicfile.New(path.Join(dir, "junk")),
		},
		{
			name:       "error",
			source:     "gs://m-lab-gcs-downloader-mlab-testing/t1/okay.txt",
			atomicFile: fakeNew(path.Join(dir, "junk"), fmt.Errorf("fake temp file create failed")),
		},
		{
			name:   "error-log",
			source: "gs://this-bucket-does-not-exist/okay.txt",
		},
	}

	for _, tt := range tests {
		once = true // Only run once.
		source = tt.source
		atomicfileNew = func(name string) atomicfile.File { return tt.atomicFile }
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

type fakeAtomicfile struct {
	name        string
	createError error
}

func fakeNew(name string, createError error) atomicfile.File {
	return &fakeAtomicfile{name: name, createError: createError}
}

func (f *fakeAtomicfile) Create(perm os.FileMode) (*safefile.File, error) {
	return &safefile.File{}, f.createError
}
func (f *fakeAtomicfile) Commit() error {
	return nil
}
func (f *fakeAtomicfile) Name() string {
	return f.name
}
