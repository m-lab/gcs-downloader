// Package atomicfile defines file operations for safely writing files.
package atomicfile

import (
	"context"
	"io"
	"log"
	"os"
	"path"

	"github.com/dchest/safefile"
)

// A File abstracts operations on a temporary file.
type File interface {
	// Name returns the file name.
	Name() string
	// Create allocates a new safefile.File that writes to a temporary file.
	Create(perm os.FileMode) (*safefile.File, error)
	// Commit closes the underlying temporary file and renames it to the File.Name().
	Commit() error
}

type tempFile struct {
	f    *safefile.File
	name string
}

// New creates a new temporary File.
func New(filename string) File {
	return &tempFile{name: filename}
}

// Name returns the temporary file name.
func (t *tempFile) Name() string {
	return t.name
}

// Create allocates a new safefile.File that writes to a temporary file.
func (t *tempFile) Create(perm os.FileMode) (*safefile.File, error) {
	var err error
	t.f, err = safefile.Create(t.name, perm)
	return t.f, err
}

// Commit closes the underlying temporary file and writes to the File.Name().
func (t *tempFile) Commit() error {
	return t.f.Commit()
}

// Object abstracts a data source.
type Object interface {
	// Name returns the Object name.
	ObjectName() string
	// Copy writes the Object data to the given writer.
	Copy(ctx context.Context, w io.Writer) error
}

// SaveFile downloads object data and saves it to the given temp file.
func SaveFile(ctx context.Context, o Object, t File) error {
	// Make dirs locally for any remaining directory path from source name.
	err := os.MkdirAll(path.Dir(t.Name()), os.ModePerm)
	if err != nil {
		log.Printf("Failed to make local directory: %s", err)
		return err
	}

	log.Println("Writing:", t.Name())
	f, err := t.Create(0644)
	if err != nil {
		log.Printf("Failed to write file %q: %s", t.Name(), err)
		return err
	}
	defer f.Close()

	// Download the object data.
	err = o.Copy(ctx, f)
	if err != nil {
		log.Printf("Failed to download object %q: %s", o.ObjectName(), err)
		return err
	}

	// Commit the temporary file to the intended output name.
	err = t.Commit()
	if err != nil {
		log.Printf("Failed to commit temporary file %q: %s", t.Name(), err)
		return err
	}
	return nil
}
