package main

import (
	"net/http"

	"github.com/areknoster/gopherdir/pkg/gopherdir"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// api manages
type api struct {
	logger      *zap.Logger
	fileManager gopherdir.FileManager
	uiRenderer  gopherdir.UIRenderer
}

func newApi(logger *zap.Logger, fileManager gopherdir.FileManager, uiRenderer gopherdir.UIRenderer) *api {
	return &api{logger: logger, fileManager: fileManager, uiRenderer: uiRenderer}
}

func (a *api) mount(router chi.Router) {
	router.Get("/", a.handleGetUI)
}

func (a *api) handleGetUI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	files, err := a.fileManager.GetFiles(ctx)
	if err != nil{
		http.Error(w, "could error fetching driectory contents", http.StatusInternalServerError)
		return
	}
	err = a.uiRenderer.Render(ctx,  files, w)
	if err != nil {
		a.logger.Info("could not render UI", zap.Error(err))
	}
}
