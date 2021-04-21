package gopherdir

import (
	"context"
	"io"
)

// UIRenderer is responsible for rendering UI for the application
type UIRenderer interface {
	Render(ctx context.Context, files []File, w io.Writer) error
}
