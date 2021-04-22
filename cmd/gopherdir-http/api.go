package main

import (
	"net/http"

	"github.com/areknoster/gopherdir/pkg/gopherdir"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type apiConfig struct {
	MaxFileBytes int64 `default:"1048576" split_words:"true" description:"max file size in bytes"`
}

// api manages
type api struct {
	logger      *zap.Logger
	fileManager gopherdir.FileManager
	uiRenderer  gopherdir.UIRenderer
	maxFileSize int64
}

func newApi(logger *zap.Logger, fileManager gopherdir.FileManager, uiRenderer gopherdir.UIRenderer, cfg apiConfig) *api {
	return &api{
		logger:      logger,
		fileManager: fileManager,
		uiRenderer:  uiRenderer,
		maxFileSize: cfg.MaxFileBytes,
	}
}

const FileNameParam = "fileName"

func (a *api) mount(router chi.Router) {
	router.Get("/", a.handleGetUI)
	router.Route("/file", func(r chi.Router) {
		r.Post("/", a.handleUploadFile)
		r.Get("/*", a.handleDownloadFile)
	})
}

func (a *api) handleGetUI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	files, err := a.fileManager.GetFileNames(ctx)
	if err != nil {
		a.logger.Info("could not fetch directory contents", zap.Error(err))
		http.Error(w, "could not fetch directory contents", http.StatusInternalServerError)
		return
	}
	err = a.uiRenderer.Render(ctx, files, w)
	if err != nil {
		a.logger.Info("could not render UI", zap.Error(err))
		http.Error(w, "could not fetch directory contents", http.StatusInternalServerError)
	}
}

func (a *api) handleUploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(a.maxFileSize); err != nil {
		a.logger.Info("could not parse request", zap.Error(err))
		http.Error(w, "could not parse request", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("upload_file")
	if err != nil {
		a.logger.Info("could not retrieve file from form", zap.Error(err))
		http.Error(w, "could not retrieve file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	if err := a.fileManager.CreateFile(r.Context(), file, header.Filename); err != nil {
		a.logger.Info("could create file", zap.Error(err))
		http.Error(w, "could create file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<html>Upload successful! <a href="/">Return to files page</a></html>`))
}

func (a *api) handleDownloadFile(w http.ResponseWriter, r *http.Request) {
	fileServer := http.FileServer(http.FS(a.fileManager.GetFileSystem()))
	http.StripPrefix("/file/", fileServer).ServeHTTP(w, r)
}
