package static

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestCollectCopiesFiles(t *testing.T) {
	src := fstest.MapFS{
		"robots.txt":       {Data: []byte("robots")},
		"css/main.css":     {Data: []byte("body{}")},
		"images/a/b/x.png": {Data: []byte("png")},
	}
	dest := t.TempDir()

	n, err := Collect(src, dest)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("expected 3 files copied, got %d", n)
	}

	for path, f := range src {
		got, err := os.ReadFile(filepath.Join(dest, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(f.Data) {
			t.Errorf("%s: expected %q, got %q", path, f.Data, got)
		}
	}

	n, err = Collect(src, dest)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected 0 files copied on second run, got %d", n)
	}

	entries, err := os.ReadDir(filepath.Join(dest, "css"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected no leftover temp files, got %d entries", len(entries))
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCollectOverwritesChangedFiles(t *testing.T) {
	src := fstest.MapFS{
		"css/main.css": {Data: []byte("new")},
		"js/app.js":    {Data: []byte("app")},
		"js/same.js":   {Data: []byte("same")},
	}
	dest := t.TempDir()

	changed := filepath.Join(dest, "css", "main.css")
	writeFile(t, changed, "old")
	writeFile(t, filepath.Join(dest, "js", "same.js"), "same")

	n, err := Collect(src, dest)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 files copied, got %d", n)
	}

	got, err := os.ReadFile(changed)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "new" {
		t.Errorf("changed file was not overwritten: got %q", got)
	}
}

func TestCollectOverwritesSameSizeDifferentContents(t *testing.T) {
	src := fstest.MapFS{
		"app.js": {Data: []byte("bbbb")},
	}
	dest := t.TempDir()
	writeFile(t, filepath.Join(dest, "app.js"), "aaaa")

	n, err := Collect(src, dest)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("expected 1 file copied, got %d", n)
	}
}

func TestCollectSkipsIdenticalFiles(t *testing.T) {
	src := fstest.MapFS{
		"app.js": {Data: []byte("same")},
	}
	dest := t.TempDir()
	writeFile(t, filepath.Join(dest, "app.js"), "same")

	n, err := Collect(src, dest)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected 0 files copied, got %d", n)
	}
}
