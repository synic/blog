package internal

import (
	"io/fs"
)

func MustSub(fsys fs.FS, path string) fs.FS {
	st, err := fs.Sub(fsys, path)

	if err != nil {
		panic(err)
	}

	return st
}
