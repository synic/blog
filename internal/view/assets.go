package view

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func BundleStaticAssets(filesystem fs.FS, files ...string) ([]byte, error) {
	var content strings.Builder
	for _, file := range files {
		data, err := fs.ReadFile(filesystem, file)
		wrap := "script"

		if filepath.Ext(file) == ".css" {
			wrap = "style"
		}

		if err != nil {
			return []byte{}, err
		}

		content.WriteString(fmt.Sprintf("<%s hx-preserve=\"true\">/* %s */\n%s</%s>", wrap, file, data, wrap))
	}

	return []byte(content.String()), nil
}
