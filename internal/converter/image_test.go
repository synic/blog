package converter

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformImagesNoImageTag(t *testing.T) {
	input := `<p>Hello</p><img src="/a.jpg" alt="A"/>`
	out, err := TransformImages(input, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestTransformImagesBasic(t *testing.T) {
	input := `<p>Before</p><x-image src="/static/foo.jpg" alt="My photo"/><p>After</p>`
	out, err := TransformImages(input, "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `class="photo lightbox-img"`)
	assert.Contains(t, out, `src="/static/foo.jpg"`)
	assert.Contains(t, out, `alt="My photo"`)
	assert.Contains(t, out, `loading="lazy"`)
	assert.NotContains(t, out, "<x-image")
	assert.Contains(t, out, "<p>Before</p>")
	assert.Contains(t, out, "<p>After</p>")
}

func TestTransformImagesManualDimensions(t *testing.T) {
	input := `<x-image src="/static/foo.jpg" alt="X" width="400" height="300"/>`
	out, err := TransformImages(input, "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `width="400"`)
	assert.Contains(t, out, `height="300"`)
}

func TestTransformImagesAutoDimensions(t *testing.T) {
	dir := t.TempDir()
	imgDir := filepath.Join(dir, "images")
	os.MkdirAll(imgDir, 0o755)

	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	for y := range 600 {
		for x := range 800 {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	f, err := os.Create(filepath.Join(imgDir, "test.jpg"))
	assert.NoError(t, err)
	assert.NoError(t, jpeg.Encode(f, img, nil))
	f.Close()

	input := `<x-image src="/static/images/test.jpg" alt="Test"/>`
	out, err := TransformImages(input, dir, nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `width="800"`)
	assert.Contains(t, out, `height="600"`)
}

func TestTransformImagesPartialManualWidth(t *testing.T) {
	dir := t.TempDir()
	imgDir := filepath.Join(dir, "images")
	os.MkdirAll(imgDir, 0o755)

	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	f, err := os.Create(filepath.Join(imgDir, "test.jpg"))
	assert.NoError(t, err)
	assert.NoError(t, jpeg.Encode(f, img, nil))
	f.Close()

	input := `<x-image src="/static/images/test.jpg" alt="Test" width="400"/>`
	out, err := TransformImages(input, dir, nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `width="400"`)
	assert.Contains(t, out, `height="300"`)
}

func TestTransformImagesPartialManualHeight(t *testing.T) {
	dir := t.TempDir()
	imgDir := filepath.Join(dir, "images")
	os.MkdirAll(imgDir, 0o755)

	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	f, err := os.Create(filepath.Join(imgDir, "test.jpg"))
	assert.NoError(t, err)
	assert.NoError(t, jpeg.Encode(f, img, nil))
	f.Close()

	input := `<x-image src="/static/images/test.jpg" alt="Test" height="300"/>`
	out, err := TransformImages(input, dir, nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `width="400"`)
	assert.Contains(t, out, `height="300"`)
}

func TestTransformImagesMissingFile(t *testing.T) {
	input := `<x-image src="/static/images/nonexistent.jpg" alt="X"/>`
	out, err := TransformImages(input, "/tmp/empty", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `src="/static/images/nonexistent.jpg"`)
	assert.NotContains(t, out, "width=")
	assert.NotContains(t, out, "height=")
}

func TestTransformImagesPreservesSurrounding(t *testing.T) {
	input := `<h2>Title</h2><x-image src="/a.jpg" alt="A"/><p>End</p>`
	out, err := TransformImages(input, "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, "<h2>Title</h2>")
	assert.Contains(t, out, "<p>End</p>")
	assert.Contains(t, out, `class="photo lightbox-img"`)
}
