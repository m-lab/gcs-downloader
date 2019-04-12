package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/dchest/safefile"
	"github.com/m-lab/gcs-downloader/atomicfile"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func Test_main(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		tempFile atomicfile.File
	}{
		{
			name:     "ok",
			source:   "gs://m-lab-gcs-downloader-mlab-testing/t1/okay.txt",
			tempFile: atomicfile.New("/tmp/junk"),
		},
		{
			name:     "error",
			source:   "gs://m-lab-gcs-downloader-mlab-testing/t1/okay.txt",
			tempFile: fakeNew("/tmp/junk", fmt.Errorf("fake temp file create failed")),
		},
	}

	for _, tt := range tests {
		once = true // Only run once.
		source = tt.source
		tempNew = func(name string) atomicfile.File { return tt.tempFile }
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

type fakeTemp struct {
	name        string
	createError error
}

func fakeNew(name string, createError error) atomicfile.File {
	return &fakeTemp{name: name, createError: createError}
}

func (f *fakeTemp) Create(perm os.FileMode) (*safefile.File, error) {
	return &safefile.File{}, f.createError
}
func (f *fakeTemp) Commit() error {
	return nil
}
func (f *fakeTemp) Name() string {
	return f.name
}
