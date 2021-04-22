package gopherdir

import (
	"context"
	"io"
	"io/fs"
)

type File struct {
	Name string `require:"true"`
}

type FileManager interface {
	GetFileNames(ctx context.Context) ([]File, error)
	CreateFile(ctx context.Context, file io.Reader, filename string) error
	GetFileSystem() fs.FS
}
