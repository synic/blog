//go:build mage

package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	_ "golang.org/x/image/webp"

	"github.com/synic/blog/internal/article"
	"github.com/synic/blog/internal/config"
	"github.com/synic/blog/internal/converter"
	"github.com/synic/blog/internal/model"
)

func init() {
	godotenv.Load()
}

var (
	Default        = Dev
	packageName    = "github.com/synic/blog"
	syntaxCssTheme = "nord-darker"

	// paths
	binPath         = "bin"
	debugPath       = path.Join(binPath, "blog-debug")
	releasePath     = path.Join(binPath, "blog-release")
	cssOutPath      = "./static/css"
	tailwinInFile   = "./assets/css/main.css"
	tailwindOutFile = cssOutPath + "/main.min.css"
	syntaxCssFile   = "./assets/css/syntax.css"
	articlesInPath  = "./assets/articles"
	articlesOutPath = "./static/articles"
	configPath      = packageName + "/internal/config"
	imagesInPath    = "./assets/images"
	imagesOutPath   = "./static/images"

	migrationsPath = "./migrations"

	// commands
	buildCmd    = sh.RunCmd("go", "build")
	tailwindCmd = sh.RunCmd("node_modules/.bin/tailwindcss")
	minifyCmd   = sh.RunCmd("node_modules/.bin/css-minify")
	terserCmd   = sh.RunCmd("node_modules/.bin/terser")
	templCmd    = sh.RunCmd("go", "tool", "templ")
	airCmd      = sh.RunCmd("go", "tool", "air")
	gooseCmd    = sh.RunCmd("go", "tool", "goose")
	sqlcCmd     = sh.RunCmd("go", "tool", "sqlc")
)

func Dev() error {
	mg.Deps(Deps.Dev)

	return airCmd()
}

type Deps mg.Namespace

func (Deps) Dev() error {
	if _, err := os.Stat("./node_modules/.bin/tailwindcss"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err = sh.RunV("npm", "i"); err != nil {
			return err
		}
	}

	return nil
}

type Build mg.Namespace

func (Build) Dev() error {
	mg.Deps(Codegen)
	mg.Deps(Articles.Convert, Images{}.Build)

	return buildCmd("-tags", "debug",
		fmt.Sprintf("-ldflags=-X %s.DebugFlag=true", configPath),
		"-o", debugPath,
		".",
	)
}

func (Build) Release() error {
	mg.Deps(Check)
	mg.Deps(Articles.convertWithGit, Images{}.Build)
	return sh.RunWithV(
		map[string]string{"CGO_ENABLED": "0"},
		"go", "build",
		"-tags", "release",
		fmt.Sprintf("-ldflags=-s -w -X %s.BuildTime=%d", configPath, time.Now().Unix()),
		"-o", releasePath,
		".",
	)
}

func MinifyJs() error {
	fmt.Println("Minifying app.js...")
	return terserCmd("./assets/js/app.js", "-o", "./static/js/app.min.js", "--compress", "--mangle")
}

func Codegen() error {
	mg.Deps(Deps.Dev)

	if err := Sqlc(); err != nil {
		return err
	}

	if err := MinifyJs(); err != nil {
		return err
	}

	err := templCmd("generate")

	if err != nil {
		return err
	}

	return maybeRunTailwind()
}

type Articles mg.Namespace

func (Articles) Convert() error {
	return convertArticles(false, false)
}

func (Articles) convertWithGit() error {
	return convertArticles(false, true)
}

func (Articles) Reconvert() error {
	return convertArticles(true, false)
}

func (Articles) Create() error {
	var title, tags string
	var createImageDir = false

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Title: ")
	if scanner.Scan() {
		title = scanner.Text()
	}

	fmt.Print("Tags: ")
	if scanner.Scan() {
		tags = scanner.Text()
	}

	if title == "" {
		fmt.Println("Title is required.")
		os.Exit(1)
	}

	fmt.Print("Create image directory [y/N]? ")

	if scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "y" {
			createImageDir = true
		}
	}

	payload := model.ArticleCreatePayload{
		Title:       title,
		Tags:        tags,
		PublishedAt: time.Now(),
	}

	fn, err := article.CreateArticle(payload, createImageDir)

	if err != nil {
		return err
	}

	return sh.RunV("nvim", fn, "-c", "/<!-- summary -->", "-c", "normal! j0", "-c", "startinsert")
}

type Images mg.Namespace

func (Images) Import(src, destDir string) error {
	name := strings.TrimSuffix(filepath.Base(src), filepath.Ext(src)) + ".webp"
	outDir := filepath.Join(imagesInPath, destDir)
	dst := filepath.Join(outDir, name)

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	width, err := imageWidth(src)
	if err != nil {
		return err
	}

	args := []string{"-q", "82", "-m", "6"}
	if width > 1600 {
		args = append(args, "-resize", "1600", "0")
		fmt.Printf("Importing %s -> %s (%dpx -> 1600px)\n", src, dst, width)
	} else {
		fmt.Printf("Importing %s -> %s (%dpx)\n", src, dst, width)
	}
	args = append(args, src, "-o", dst)

	return sh.RunV("cwebp", args...)
}

func (Images) Build() error {
	type variant struct {
		suffix   string
		maxWidth int
	}

	sizeVariants := map[string]variant{
		"xl": {"", 0},
		"lg": {"-lg", 1000},
		"md": {"-md", 640},
		"sm": {"-sm", 480},
	}

	xImageSizes, err := collectXImageSizes()
	if err != nil {
		return fmt.Errorf("error collecting x-image sizes: %w", err)
	}

	built := 0
	copied := 0
	skipped := 0

	err = filepath.WalkDir(imagesInPath, func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		ext := strings.ToLower(filepath.Ext(srcPath))
		if ext != ".webp" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			return nil
		}

		rel, _ := filepath.Rel(imagesInPath, srcPath)
		srcInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		sizes, isXImage := xImageSizes[rel]
		if !isXImage {
			dstPath := filepath.Join(imagesOutPath, rel)
			dstInfo, dstErr := os.Stat(dstPath)
			if dstErr == nil && dstInfo.ModTime().After(srcInfo.ModTime()) {
				skipped++
				return nil
			}
			if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
				return err
			}
			if err := copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("error copying %s: %w", rel, err)
			}
			copied++
			return nil
		}

		base := strings.TrimSuffix(rel, filepath.Ext(rel))

		srcWidth, err := imageWidth(srcPath)
		if err != nil {
			fmt.Printf("⚠️  skipping %s: %v\n", rel, err)
			return nil
		}

		for _, sizeName := range sizes {
			v, ok := sizeVariants[sizeName]
			if !ok {
				fmt.Printf("⚠️  unknown size %q for %s, skipping\n", sizeName, rel)
				continue
			}

			if v.maxWidth > 0 && srcWidth <= v.maxWidth {
				continue
			}

			dstPath := filepath.Join(imagesOutPath, base+v.suffix+".webp")

			dstInfo, dstErr := os.Stat(dstPath)
			if dstErr == nil && dstInfo.ModTime().After(srcInfo.ModTime()) {
				skipped++
				continue
			}

			if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
				return err
			}

			args := []string{"-q", "82", "-m", "6"}
			if v.maxWidth > 0 {
				args = append(args, "-resize", fmt.Sprintf("%d", v.maxWidth), "0")
			}
			args = append(args, srcPath, "-o", dstPath)

			if err := sh.Run("cwebp", args...); err != nil {
				return fmt.Errorf("error converting %s: %w", rel, err)
			}

			built++
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("🖼️  Images: %d built, %d copied, %d up-to-date\n", built, copied, skipped)
	return nil
}

func imageWidth(src string) (int, error) {
	f, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("error opening image %s: %w", src, err)
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, fmt.Errorf("error decoding image %s: %w", src, err)
	}

	return cfg.Width, nil
}

var (
	xImageTagRe  = regexp.MustCompile(`<x-image\b([^>]*)/?>`)
	xImageAttrRe = regexp.MustCompile(`(\w[-\w]*)="([^"]*)"`)
)

func collectXImageSizes() (map[string][]string, error) {
	allSizes := []string{"xl", "lg", "md", "sm"}
	result := make(map[string][]string)

	addSizes := func(src, sizesStr string) {
		if src == "" || !strings.HasPrefix(src, "/static/images/") {
			return
		}
		if sizesStr == "original" {
			return
		}
		rel := strings.TrimPrefix(src, "/static/images/")
		var sizes []string
		if sizesStr != "" {
			for _, s := range strings.Split(sizesStr, ",") {
				if s = strings.TrimSpace(s); s != "" {
					sizes = append(sizes, s)
				}
			}
		} else {
			sizes = allSizes
		}
		if existing, ok := result[rel]; ok {
			sizeSet := make(map[string]bool)
			for _, s := range existing {
				sizeSet[s] = true
			}
			for _, s := range sizes {
				if !sizeSet[s] {
					existing = append(existing, s)
					sizeSet[s] = true
				}
			}
			result[rel] = existing
		} else {
			result[rel] = sizes
		}
	}

	err := filepath.WalkDir(articlesInPath, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Ext(p) != ".md" {
			return nil
		}

		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		content := string(data)

		for _, match := range xImageTagRe.FindAllString(content, -1) {
			attrs := parseAttrs(match)
			addSizes(attrs["src"], attrs["sizes"])
		}

		return nil
	})

	return result, err
}

func parseAttrs(tag string) map[string]string {
	attrs := make(map[string]string)
	for _, m := range xImageAttrRe.FindAllStringSubmatch(tag, -1) {
		attrs[m[1]] = m[2]
	}
	return attrs
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

type DB mg.Namespace

func databaseUrl() string {
	conf := config.Load()
	return conf.DatabaseURL
}

func (DB) Migrate() error {
	return gooseCmd("-dir", migrationsPath, "sqlite3", databaseUrl(), "up")
}

func (DB) Rollback() error {
	return gooseCmd("-dir", migrationsPath, "sqlite3", databaseUrl(), "down")
}

func (DB) Status() error {
	return gooseCmd("-dir", migrationsPath, "sqlite3", databaseUrl(), "status")
}

func Sqlc() error {
	return sqlcCmd("generate")
}

func Pygmentize() error {
	data, err := sh.Output("pygmentize", "-S", syntaxCssTheme, "-f", "html", "-a", ".chroma")

	if err != nil {
		return err
	}

	os.Remove(syntaxCssFile)
	f, err := os.Create(syntaxCssFile)

	if err != nil {
		return err
	}

	f.WriteString(data)
	f.Close()

	return minifyCmd("-f", syntaxCssFile, "--output", cssOutPath)
}

func Vet() error {
	return sh.RunV("go", "vet", "./...")
}

func Test() error {
	mg.Deps(Vet)
	return sh.RunV("go", "test", "-race", "./...")
}

func Check() error {
	if err := templCmd("fmt", "internal/view"); err != nil {
		return err
	}

	if err := sh.RunV("go", "fmt", "./..."); err != nil {
		return err
	}

	if err := Vet(); err != nil {
		return err
	}

	return sh.RunV("go", "test", "-race", "./...")
}

func Clean() error {
	files, err := os.ReadDir(binPath)

	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name() == ".gitkeep" {
			continue
		}

		err := sh.Rm(path.Join(binPath, file.Name()))

		if err != nil {
			return err
		}
	}

	return nil
}

func Container() error {
	return sh.RunV("docker", "build", "-t", "blog", ".")
}

func lastModTime(paths ...string) (*time.Time, error) {
	var lastModTime *time.Time = nil

	for _, path := range paths {
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			fileInfo, err := os.Stat(path)

			if err != nil {
				return err
			}

			modTime := fileInfo.ModTime()

			if lastModTime == nil || modTime.After(*lastModTime) {
				lastModTime = &modTime
			}

			return nil
		})
		if err != nil {
			return lastModTime, err
		}
	}

	return lastModTime, nil
}

func maybeRunTailwind() error {
	var outModTime *time.Time = nil

	outInfo, err := os.Stat(tailwindOutFile)
	if err == nil {
		modTime := outInfo.ModTime()
		outModTime = &modTime
	} else if !os.IsNotExist(err) {
		return err
	}

	lastInputModTime, err := lastModTime("./internal/view", articlesInPath)
	if err != nil {
		return err
	}

	if lastInputModTime == nil || outModTime == nil || lastInputModTime.After(*outModTime) {
		fmt.Println("View or CSS changes detected, running tailwind...")
		if err := os.MkdirAll(filepath.Dir(tailwindOutFile), 0o755); err != nil {
			return err
		}
		err := tailwindCmd("-i", tailwinInFile, "-o", tailwindOutFile, "--minify")
		if err != nil {
			return err
		}

		return sh.Run("touch", tailwindOutFile)
	}

	return nil
}

func convertArticles(reconvert, useGit bool) error {
	res, err := converter.Convert(articlesInPath, articlesOutPath, reconvert, useGit)

	if err != nil {
		return err
	}

	log.New(os.Stderr, "", log.LstdFlags).Print(res.String())
	return nil
}
