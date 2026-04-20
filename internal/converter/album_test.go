package converter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformAlbumsNoAlbums(t *testing.T) {
	input := `<p>Hello world</p><img src="/a.jpg" alt="a"/>`
	out, err := TransformAlbums(input, "b", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, input, out)
}

func TestTransformAlbumsSingleWithCaptions(t *testing.T) {
	input := `<x-album>
<x-image src="/a.jpg" alt="First"/>
<x-image src="/b.jpg" alt="Second"/>
<x-image src="/c.jpg" alt="Third"/>
</x-album>`
	out, err := TransformAlbums(input, "b", "", nil)
	assert.NoError(t, err)

	assert.Contains(t, out, `class="album"`)
	assert.Contains(t, out, `data-album-id="album-b-0"`)
	assert.Contains(t, out, `class="album-frame"`)
	assert.Contains(t, out, `class="album-nav album-nav-prev"`)
	assert.Contains(t, out, `class="album-nav album-nav-next"`)
	assert.Contains(t, out, `class="album-img`)
	assert.Contains(t, out, `src="/a.jpg"`)
	assert.Contains(t, out, `alt="First"`)
	assert.Contains(t, out, `<div class="album-caption" aria-live="polite">First</div>`)
	assert.Contains(t, out, `class="album-dots"`)
	assert.Contains(t, out, `role="region"`)
	assert.Contains(t, out, `aria-label="Image album"`)
	assert.Contains(t, out, `aria-label="Previous image"`)
	assert.Contains(t, out, `aria-label="Next image"`)
	assert.Contains(t, out, `class="album-dot is-active"`)
	assert.Contains(t, out, `data-index="2"`)
	assert.NotContains(t, out, "<x-album")
}

func TestTransformAlbumsMissingAlt(t *testing.T) {
	input := `<x-album><x-image src="/a.jpg"/></x-album>`
	out, err := TransformAlbums(input, "b", "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `<div class="album-caption" aria-live="polite"></div>`)
	assert.Contains(t, out, `src="/a.jpg"`)
}

func TestTransformAlbumsEmpty(t *testing.T) {
	input := `<p>Before</p><x-album></x-album><p>After</p>`
	out, err := TransformAlbums(input, "b", "", nil)
	assert.NoError(t, err)
	assert.NotContains(t, out, "album")
	assert.Contains(t, out, "<p>Before</p>")
	assert.Contains(t, out, "<p>After</p>")
}

func TestTransformAlbumsMultiple(t *testing.T) {
	input := `<x-album><x-image src="/a.jpg" alt="A"/></x-album><p>between</p><x-album><x-image src="/b.jpg" alt="B"/></x-album>`
	out, err := TransformAlbums(input, "b", "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `data-album-id="album-b-0"`)
	assert.Contains(t, out, `data-album-id="album-b-1"`)
	assert.Equal(t, 2, strings.Count(out, `class="album"`))
}

func TestTransformAlbumsCustomPrefix(t *testing.T) {
	input := `<x-album><x-image src="/a.jpg" alt="A"/></x-album>`
	out, err := TransformAlbums(input, "summary", "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `data-album-id="album-summary-0"`)
}

func TestTransformAlbumsImgSkippedWithoutSrc(t *testing.T) {
	input := `<x-album><x-image alt="no src"/><x-image src="/a.jpg" alt="A"/></x-album>`
	out, err := TransformAlbums(input, "b", "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, `src="/a.jpg"`)
	assert.Equal(t, 2, strings.Count(out, "album-dot"))
}

func TestTransformAlbumsPreservesSurroundingContent(t *testing.T) {
	input := `<h2>Title</h2><p>Intro</p><x-album><x-image src="/a.jpg" alt="A"/></x-album><p>Outro</p>`
	out, err := TransformAlbums(input, "b", "", nil)
	assert.NoError(t, err)
	assert.Contains(t, out, "<h2>Title</h2>")
	assert.Contains(t, out, "<p>Intro</p>")
	assert.Contains(t, out, "<p>Outro</p>")
	assert.Contains(t, out, `class="album"`)
}
