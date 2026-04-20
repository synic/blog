package converter

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "golang.org/x/image/webp"
)

var imageTagRe = regexp.MustCompile(`<x-image\b([^>]*)/?>`)
var imageAttrRe = regexp.MustCompile(`(\w+)="([^"]*)"`)

func TransformImages(htmlStr, staticDir string, firstSeen *bool) (string, error) {
	if !strings.Contains(htmlStr, "<x-image") {
		return htmlStr, nil
	}

	result := imageTagRe.ReplaceAllStringFunc(htmlStr, func(match string) string {
		attrs := parseImageAttrs(match)
		src := attrs["src"]
		alt := attrs["alt"]
		widthStr := attrs["width"]
		heightStr := attrs["height"]

		isFirst := false
		if firstSeen != nil && !*firstSeen {
			isFirst = true
			*firstSeen = true
		}

		if widthStr == "" || heightStr == "" {
			w, h := resolveImageDimensions(src, staticDir)
			if w > 0 && h > 0 {
				if widthStr == "" && heightStr == "" {
					widthStr = strconv.Itoa(w)
					heightStr = strconv.Itoa(h)
				} else if widthStr == "" {
					manual, err := strconv.Atoi(heightStr)
					if err == nil && h > 0 {
						widthStr = strconv.Itoa(int(math.Round(float64(w) * float64(manual) / float64(h))))
					}
				} else {
					manual, err := strconv.Atoi(widthStr)
					if err == nil && w > 0 {
						heightStr = strconv.Itoa(int(math.Round(float64(h) * float64(manual) / float64(w))))
					}
				}
			}
		}

		var b strings.Builder
		loading := "lazy"
		fetchpriority := ""
		decoding := "async"
		if isFirst {
			loading = "eager"
			fetchpriority = ` fetchpriority="high"`
			decoding = "sync"
		}

		b.WriteString(`<picture>`)
		if srcsetWebp := buildSrcset(src, staticDir, ".webp"); srcsetWebp != "" {
			fmt.Fprintf(&b, `<source type="image/webp" srcset="%s" sizes="(max-width: 640px) 100vw, 75vw"/>`, srcsetWebp)
		}

		b.WriteString(`<img class="photo lightbox-img"`)
		fmt.Fprintf(&b, ` src="%s"`, src)
		fmt.Fprintf(&b, ` alt="%s"`, alt)
		fmt.Fprintf(&b, ` loading="%s"`, loading)
		fmt.Fprintf(&b, ` decoding="%s"`, decoding)
		b.WriteString(fetchpriority)

		if widthStr != "" {
			fmt.Fprintf(&b, ` width="%s"`, widthStr)
		}
		if heightStr != "" {
			fmt.Fprintf(&b, ` height="%s"`, heightStr)
		}
		if srcset := buildSrcset(src, staticDir, ""); srcset != "" {
			fmt.Fprintf(&b, ` srcset="%s"`, srcset)
			b.WriteString(` sizes="(max-width: 640px) 100vw, 75vw"`)
		}
		b.WriteString("/></picture>")
		return b.String()
	})

	return result, nil
}

func parseImageAttrs(tag string) map[string]string {
	attrs := make(map[string]string)
	matches := imageAttrRe.FindAllStringSubmatch(tag, -1)
	for _, m := range matches {
		attrs[m[1]] = m[2]
	}
	return attrs
}

func resolveImageDimensions(src, staticDir string) (int, int) {
	if staticDir == "" || !strings.HasPrefix(src, "/static/") {
		return 0, 0
	}

	rel := strings.TrimPrefix(src, "/static/")
	path := filepath.Join(staticDir, rel)

	f, err := os.Open(path)
	if err != nil {
		return 0, 0
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0
	}

	return cfg.Width, cfg.Height
}

type srcsetEntry struct {
	url   string
	width int
}

func buildSrcset(src, staticDir, forceExt string) string {
	if staticDir == "" || !strings.HasPrefix(src, "/static/") {
		return ""
	}

	ext := filepath.Ext(src)
	base := strings.TrimSuffix(src, ext)

	if forceExt != "" {
		ext = forceExt
	}

	suffixes := []string{"-sm", "-md", "-lg", ""}
	var entries []srcsetEntry

	for _, suffix := range suffixes {
		var candidateURL string
		if suffix == "" {
			candidateURL = base + ext
		} else {
			candidateURL = base + suffix + ext
		}

		w, _ := resolveImageDimensions(candidateURL, staticDir)
		if w > 0 {
			entries = append(entries, srcsetEntry{url: candidateURL, width: w})
		}
	}

	if len(entries) == 0 {
		return ""
	}

	if len(entries) <= 1 && forceExt == "" {
		return ""
	}

	var parts []string
	for _, e := range entries {
		parts = append(parts, fmt.Sprintf("%s %dw", e.url, e.width))
	}

	return strings.Join(parts, ", ")
}
