package temp_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/dchest/safefile"
	"github.com/m-lab/gcs-downloader/temp"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestNew(t *testing.T) {
	tmp := temp.New("test1")
	_, err := tmp.Create(os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = tmp.Commit()
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(tmp.Name())
}

type fakeTemp struct {
	name        string
	createError error
	commitError error
}

func fakeFile(name string, createError, commitError error) temp.File {
	return &fakeTemp{name: name, createError: createError, commitError: commitError}
}

func (f *fakeTemp) Create(perm os.FileMode) (*safefile.File, error) {
	if f.createError != nil {
		return nil, f.createError
	}
	return safefile.Create(f.name, perm)
}
func (f *fakeTemp) Commit() error {
	return f.commitError
}
func (f *fakeTemp) Name() string {
	return f.name
}

type fakeObj struct {
	name      string
	buf       *bytes.Buffer
	copyError error
}

func (f *fakeObj) ObjectName() string {
	return f.name
}
func (f *fakeObj) LocalName() string {
	return f.name
}
func (f *fakeObj) Copy(ctx context.Context, w io.Writer) error {
	if f.copyError != nil {
		return f.copyError
	}
	_, err := io.Copy(w, f.buf)
	return err
}

func TestSaveFile(t *testing.T) {
	tests := []struct {
		name    string
		o       temp.Object
		t       temp.File
		mkdir   bool
		wantErr bool
	}{
		{
			name: "okay",
			o:    &fakeObj{name: "okay", buf: bytes.NewBufferString("test")},
			t:    fakeFile("local", nil, nil),
		},
		{
			name:    "error-mkdir",
			t:       fakeFile("fakedir/local", nil, nil),
			mkdir:   true,
			wantErr: true,
		},
		{
			name:    "error-create",
			t:       fakeFile("local", fmt.Errorf("Fake create error"), nil),
			wantErr: true,
		},
		{
			name:    "error-copy",
			o:       &fakeObj{name: "okay", copyError: fmt.Errorf("Fake copy error")},
			t:       fakeFile("local", nil, nil),
			wantErr: true,
		},
		{
			name:    "error-commit",
			o:       &fakeObj{name: "okay", buf: bytes.NewBufferString("test")},
			t:       fakeFile("local", nil, fmt.Errorf("Fake commit error")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		if tt.mkdir {
			_, err := os.Create("fakedir")
			if err != nil {
				t.Fatal(err)
			}
			defer func() { os.Remove("fakedir") }()
		}

		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if err := temp.SaveFile(ctx, tt.o, tt.t); (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
