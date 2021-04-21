package ui

import (
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/areknoster/gopherdir/pkg/gopherdir"
)

// Renderer renders user content. It implements gopherdir.Renderer interface.
type Renderer struct{}

func (r Renderer) Render(ctx context.Context, files []gopherdir.File, w io.Writer) error {
	tmpl, err := template.ParseFS(fs, "hello.gohtml")
	if err != nil {
		return fmt.Errorf("ui-renderer: could not parse file: %w", err)
	}

	err = tmpl.Execute(w, &DirBrowser{Files: files})
	if err != nil {
		return fmt.Errorf("ui-renderer: could not render template: %w", err)
	}
	return nil
}

var _ gopherdir.UIRenderer = &Renderer{}
