package internal

import (
	"fmt"
	"io/fs"
	"strings"
)

func BundleStaticAssets(filesystem fs.FS, files ...string) ([]byte, error) {
	content := ""
	for _, file := range files {
		data, err := fs.ReadFile(filesystem, file)
		wrap := "script"

		if strings.HasSuffix(file, ".css") {
			wrap = "style"
		}

		if err != nil {
			return []byte{}, err
		}

		content += fmt.Sprintf("<%s>/* %s */\n%s</%s>", wrap, file, data, wrap)
	}

	return []byte(content), nil
}
