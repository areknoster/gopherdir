package files

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/areknoster/gopherdir/pkg/gopherdir"
	"go.uber.org/zap"
)

type FileManagerConfig struct {
	Root    string `required:"true"`
	BufSize int    `default:"1024"`
}

type FileManager struct {
	logger  *zap.Logger
	root    string
	bufSize int
}

func NewFileManager(logger *zap.Logger, config FileManagerConfig) (*FileManager, error) {
	// check if directory is readable, otherwise fail instantly
	_, err := os.ReadDir(config.Root)
	if err != nil {
		return nil, fmt.Errorf("file-manager: could not read directory: %w", err)
	}

	return &FileManager{
		logger:  logger,
		root:    config.Root,
		bufSize: config.BufSize,
	}, nil
}

var _ gopherdir.FileManager = &FileManager{}

func (f *FileManager) GetFileNames(ctx context.Context) ([]gopherdir.File, error) {
	entries, err := os.ReadDir(f.root)
	if err != nil {
		return nil, fmt.Errorf("file-manager: could not read directory: %w", err)
	}

	return entriesToFiles(entries), nil
}

func entriesToFiles(entries []os.DirEntry) []gopherdir.File {
	files := make([]gopherdir.File, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, gopherdir.File{
				Name: entry.Name(),
			})
		}
	}
	return files
}

func (f *FileManager) CreateFile(ctx context.Context, content io.Reader, filename string) error {
	localFile, err := os.Create(f.root + "/" + filename)
	if err != nil {
		return fmt.Errorf("could not save content to local directory: %w", err)
	}
	defer localFile.Close()

	buf := make([]byte, f.bufSize)
	if _, err = io.CopyBuffer(localFile, content, buf); err != nil {
		return fmt.Errorf("could not copy file content to local file: %w", err)
	}
	return nil
}

func (f *FileManager) GetFileSystem() fs.FS {
	return os.DirFS(f.root)
}
