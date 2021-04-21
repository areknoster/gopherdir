package files

import (
	"context"
	"fmt"
	"github.com/areknoster/gopherdir/pkg/gopherdir"
	"go.uber.org/zap"
	"os"
)

type FileManagerConfig struct{
	Root string `required:"true"`
}

type FileManager struct {
	logger *zap.Logger
	root   string
}

func NewFileManager(logger *zap.Logger, config FileManagerConfig) (*FileManager, error) {
	// check if directory is readable, otherwise fail instantly
	_, err := os.ReadDir(config.Root)
	if err != nil {
		return nil, fmt.Errorf("file-manager: could not read directory: %w", err)
	}

	return &FileManager{logger: logger, root: config.Root}, nil
}

var _ gopherdir.FileManager = &FileManager{}


func (f *FileManager) GetFiles(ctx context.Context) ([]gopherdir.File, error) {
	entries, err := os.ReadDir(f.root)
	if err != nil {
		return nil, fmt.Errorf("file-manager: could not read directory: %w", err)
	}

	return entriesToFiles(entries), nil
}

func entriesToFiles(entries []os.DirEntry) []gopherdir.File{
	files := make([]gopherdir.File, 0)
	for _, entry := range entries {
		if !entry.IsDir(){
			files = append(files, gopherdir.File{
				Name: entry.Name(),
			})
		}
	}
	return files
}