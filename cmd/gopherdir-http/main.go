package main

import (
	"fmt"
	"log"

	"github.com/areknoster/gopherdir/pkg/files"

	"github.com/areknoster/gopherdir/pkg/app"
	"github.com/areknoster/gopherdir/pkg/ui"

	"go.uber.org/zap"
)

// config defines the configuration for application. The config is set via env arguments
// 		run with --help flag for envs documentation
type config struct {
	App         app.Config              `split_words:"true"`
	FileManager files.FileManagerConfig `split_words:"true"`
	Api         apiConfig               `split_words:"true"`
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Panicf("could not initialize logger: %s", err.Error())
	}
	var cfg config
	app.LoadConfig(logger, &cfg)

	application := app.NewApp(logger, cfg.App)
	application.Build(func() error {
		fileManager, err := files.NewFileManager(application.Logger, cfg.FileManager)
		if err != nil {
			return fmt.Errorf("could not initialize file manager: %w", err)
		}

		api := newApi(logger, fileManager, ui.Renderer{}, cfg.Api)
		api.mount(application.Router)
		return nil
	})
	application.Run()
}
