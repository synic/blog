package static

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Collect(src fs.FS, dest string) (int, error) {
	copied := 0

	err := fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		target := filepath.Join(dest, filepath.FromSlash(path))

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		if !d.Type().IsRegular() {
			return nil
		}

		same, err := sameContents(src, path, target)
		if err != nil {
			return err
		}

		if same {
			return nil
		}

		if err := copyFile(src, path, target); err != nil {
			return err
		}

		copied++
		return nil
	})

	return copied, err
}

func sameContents(src fs.FS, path, target string) (bool, error) {
	info, err := os.Lstat(target)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	srcInfo, err := fs.Stat(src, path)
	if err != nil {
		return false, err
	}

	if !info.Mode().IsRegular() || info.Size() != srcInfo.Size() {
		return false, nil
	}

	a, err := fs.ReadFile(src, path)
	if err != nil {
		return false, err
	}

	b, err := os.ReadFile(target)
	if err != nil {
		return false, err
	}

	return bytes.Equal(a, b), nil
}

func copyFile(src fs.FS, path, target string) error {
	in, err := src.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp, err := os.CreateTemp(filepath.Dir(target), ".collect-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, in); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), target)
}
