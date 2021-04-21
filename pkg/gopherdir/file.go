package gopherdir

import "context"

type File struct{
    Name string `require:"true"`
}

type FileManager interface{
    GetFiles(ctx context.Context) ([]File, error)
}
