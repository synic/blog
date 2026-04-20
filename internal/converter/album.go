package converter

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type albumImage struct {
	Src    string
	Alt    string
	Srcset string
	Webp   string
	Width  int
	Height int
}

func TransformAlbums(htmlStr, idPrefix, staticDir string, firstSeen *bool) (string, error) {
	if !strings.Contains(htmlStr, "<x-album") {
		return htmlStr, nil
	}

	wrapped := "<html><body>" + htmlStr + "</body></html>"
	doc, err := html.Parse(strings.NewReader(wrapped))
	if err != nil {
		return "", fmt.Errorf("error parsing html for albums: %w", err)
	}

	body := findBody(doc)
	if body == nil {
		return htmlStr, nil
	}

	counter := 0
	transformAlbumNodes(body, idPrefix, staticDir, &counter, firstSeen)

	var buf strings.Builder
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		if err := html.Render(&buf, c); err != nil {
			return "", fmt.Errorf("error rendering transformed html: %w", err)
		}
	}

	return buf.String(), nil
}

func findBody(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.DataAtom == atom.Body {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if b := findBody(c); b != nil {
			return b
		}
	}
	return nil
}

func transformAlbumNodes(n *html.Node, idPrefix, staticDir string, counter *int, firstSeen *bool) {
	var next *html.Node
	for c := n.FirstChild; c != nil; c = next {
		next = c.NextSibling
		if c.Type == html.ElementNode && c.Data == "x-album" {
			isFirst := false
			if firstSeen != nil && !*firstSeen {
				isFirst = true
				*firstSeen = true
			}

			replacement := buildAlbumReplacement(c, idPrefix, staticDir, counter, isFirst)
			if replacement == nil {
				n.RemoveChild(c)
				continue
			}
			n.InsertBefore(replacement, c)
			n.RemoveChild(c)
			continue
		}
		transformAlbumNodes(c, idPrefix, staticDir, counter, firstSeen)
	}
}

func buildAlbumReplacement(albumNode *html.Node, idPrefix, staticDir string, counter *int, isFirst bool) *html.Node {
	images := collectAlbumImages(albumNode, staticDir)
	if len(images) == 0 {
		return nil
	}

	id := fmt.Sprintf("album-%s-%d", idPrefix, *counter)
	*counter++

	container := &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
		Attr: []html.Attribute{
			{Key: "class", Val: "album"},
			{Key: "role", Val: "region"},
			{Key: "aria-label", Val: "Image album"},
			{Key: "data-album-id", Val: id},
		},
	}

	frameW, frameH := 0, 0
	for _, im := range images {
		if im.Width > 0 && im.Height > 0 {
			if frameW == 0 || float64(im.Width)/float64(im.Height) < float64(frameW)/float64(frameH) {
				frameW, frameH = im.Width, im.Height
			}
		}
	}

	frameStyle := ""
	if frameW > 0 && frameH > 0 {
		frameStyle = fmt.Sprintf("aspect-ratio:%d/%d", frameW, frameH)
	}

	frameAttrs := []html.Attribute{{Key: "class", Val: "album-frame"}}
	if frameStyle != "" {
		frameAttrs = append(frameAttrs, html.Attribute{Key: "style", Val: frameStyle})
	}

	frame := &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
		Attr:     frameAttrs,
	}

	scroller := &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
		Attr: []html.Attribute{
			{Key: "class", Val: "album-scroller"},
		},
	}

	for i, imgData := range images {
		item := &html.Node{
			Type:     html.ElementNode,
			Data:     "div",
			DataAtom: atom.Div,
			Attr: []html.Attribute{
				{Key: "class", Val: "album-item"},
				{Key: "data-index", Val: strconv.Itoa(i)},
			},
		}

		picture := &html.Node{
			Type:     html.ElementNode,
			Data:     "picture",
			DataAtom: atom.Picture,
		}

		if imgData.Webp != "" {
			picture.AppendChild(&html.Node{
				Type:     html.ElementNode,
				Data:     "source",
				DataAtom: atom.Source,
				Attr: []html.Attribute{
					{Key: "type", Val: "image/webp"},
					{Key: "srcset", Val: imgData.Webp},
					{Key: "sizes", Val: "(max-width: 640px) 100vw, 75vw"},
				},
			})
		}

		loading := "lazy"
		decoding := "async"
		fetchpriority := ""
		if isFirst && i == 0 {
			loading = "eager"
			decoding = "sync"
			fetchpriority = "high"
		}

		imgAttrs := []html.Attribute{
			{Key: "class", Val: "album-img lightbox-img"},
			{Key: "src", Val: imgData.Src},
			{Key: "alt", Val: imgData.Alt},
			{Key: "loading", Val: loading},
			{Key: "decoding", Val: decoding},
		}
		if fetchpriority != "" {
			imgAttrs = append(imgAttrs, html.Attribute{Key: "fetchpriority", Val: fetchpriority})
		}
		if imgData.Width > 0 {
			imgAttrs = append(imgAttrs, html.Attribute{Key: "width", Val: strconv.Itoa(imgData.Width)})
		}
		if imgData.Height > 0 {
			imgAttrs = append(imgAttrs, html.Attribute{Key: "height", Val: strconv.Itoa(imgData.Height)})
		}
		if imgData.Srcset != "" {
			imgAttrs = append(imgAttrs, html.Attribute{Key: "srcset", Val: imgData.Srcset})
			imgAttrs = append(imgAttrs, html.Attribute{Key: "sizes", Val: "(max-width: 640px) 100vw, 75vw"})
		}

		img := &html.Node{
			Type:     html.ElementNode,
			Data:     "img",
			DataAtom: atom.Img,
			Attr:     imgAttrs,
		}

		picture.AppendChild(img)
		item.AppendChild(picture)
		scroller.AppendChild(item)
	}

	prev := &html.Node{
		Type:     html.ElementNode,
		Data:     "button",
		DataAtom: atom.Button,
		Attr: []html.Attribute{
			{Key: "type", Val: "button"},
			{Key: "class", Val: "album-nav album-nav-prev"},
			{Key: "aria-label", Val: "Previous image"},
			{Key: "data-action", Val: "prev"},
		},
	}
	prev.AppendChild(&html.Node{Type: html.TextNode, Data: "❮"})

	nextBtn := &html.Node{
		Type:     html.ElementNode,
		Data:     "button",
		DataAtom: atom.Button,
		Attr: []html.Attribute{
			{Key: "type", Val: "button"},
			{Key: "class", Val: "album-nav album-nav-next"},
			{Key: "aria-label", Val: "Next image"},
			{Key: "data-action", Val: "next"},
		},
	}
	nextBtn.AppendChild(&html.Node{Type: html.TextNode, Data: "❯"})

	frame.AppendChild(scroller)
	frame.AppendChild(prev)
	frame.AppendChild(nextBtn)

	caption := &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
		Attr: []html.Attribute{
			{Key: "class", Val: "album-caption"},
			{Key: "aria-live", Val: "polite"},
		},
	}
	caption.AppendChild(&html.Node{Type: html.TextNode, Data: images[0].Alt})

	dots := &html.Node{
		Type:     html.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
		Attr: []html.Attribute{
			{Key: "class", Val: "album-dots"},
		},
	}
	for i := range images {
		cls := "album-dot"
		if i == 0 {
			cls += " is-active"
		}
		dot := &html.Node{
			Type:     html.ElementNode,
			Data:     "button",
			DataAtom: atom.Button,
			Attr: []html.Attribute{
				{Key: "type", Val: "button"},
				{Key: "class", Val: cls},
				{Key: "data-index", Val: strconv.Itoa(i)},
				{Key: "aria-label", Val: fmt.Sprintf("Go to image %d", i+1)},
			},
		}
		dots.AppendChild(dot)
	}

	container.AppendChild(frame)
	container.AppendChild(caption)
	container.AppendChild(dots)

	return container
}

func collectAlbumImages(albumNode *html.Node, staticDir string) []albumImage {
	var buf strings.Builder
	for c := albumNode.FirstChild; c != nil; c = c.NextSibling {
		html.Render(&buf, c)
	}

	var images []albumImage
	for _, match := range imageTagRe.FindAllString(buf.String(), -1) {
		attrs := parseImageAttrs(match)
		src := attrs["src"]
		if src == "" {
			continue
		}
		alt := attrs["alt"]
		w, h := resolveImageDimensions(src, staticDir)
		srcset := buildSrcset(src, staticDir, "")
		webp := buildSrcset(src, staticDir, ".webp")
		images = append(images, albumImage{
			Src:    src,
			Alt:    alt,
			Srcset: srcset,
			Webp:   webp,
			Width:  w,
			Height: h,
		})
	}
	return images
}
