//go:build mage

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/synic/blog/internal/converter"
)

var (
	Default        = Dev
	packageName    = "github.com/synic/blog"
	syntaxCssTheme = "nord-darker"

	// paths
	binPath         = "bin"
	debugPath       = path.Join(binPath, "blog-debug")
	releasePath     = path.Join(binPath, "blog-release")
	cssOutPath      = "./assets/css"
	tailwinInFile   = "./internal/view/css/main.css"
	tailwindOutFile = cssOutPath + "/main.css"
	syntaxCssFile   = "./internal/view/css/syntax.css"
	articlesInPath  = "./articles"
	articlesOutPath = "./assets/articles"
	buildInfoPath   = packageName + "/internal"

	// commands
	runCmd      = sh.RunCmd("go", "run")
	buildCmd    = sh.RunCmd("go", "build")
	tailwindCmd = sh.RunCmd("node_modules/.bin/tailwindcss")
	minifyCmd   = sh.RunCmd("node_modules/.bin/css-minify")
	templCmd    = sh.RunCmd("go", "tool", "github.com/a-h/templ/cmd/templ")
	airCmd      = sh.RunCmd("go", "tool", "github.com/air-verse/air")
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
	mg.Deps(Codegen, Vet)
	mg.Deps(Articles.Convert)

	return buildCmd("-tags", "debug",
		fmt.Sprintf("-ldflags=-X %s.DebugFlag=true", buildInfoPath),
		"-o", debugPath,
		".",
	)
}

func (Build) Release() error {
	mg.Deps(Test)
	mg.Deps(Articles.Convert)
	return sh.RunWithV(
		map[string]string{"CGO_ENABLED": "0"},
		"go", "build",
		"-tags", "release",
		fmt.Sprintf("-ldflags=-s -w -X %s.BuildTime=%d", buildInfoPath, time.Now().Unix()),
		"-o", releasePath,
		".",
	)
}

func Codegen() error {
	mg.Deps(Deps.Dev)

	err := templCmd("generate")

	if err != nil {
		return err
	}

	return maybeRunTailwind()
}

type Articles mg.Namespace

func (Articles) Convert() error {
	return convertArticles(false)
}

func (Articles) Reconvert() error {
	return convertArticles(true)
}

func (Articles) Create() error {
	var title, tags string

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Title: ")
	if scanner.Scan() {
		title = scanner.Text()
	}

	fmt.Print("Tags: ")
	if scanner.Scan() {
		tags = scanner.Text()
	}

	if title == "" || tags == "" {
		fmt.Println("Title and tags are required.")
		os.Exit(1)
	}

	re := regexp.MustCompile(`[^\w ]+`)
	slug := re.ReplaceAllString(strings.ToLower(title), "")
	slug = strings.ReplaceAll(slug, " ", "-")

	fmt.Println(slug)
	now := time.Now()

	fn := fmt.Sprintf(
		"%s/%d-%02d-%02d_%s.md",
		articlesInPath,
		now.Year(),
		now.Month(),
		now.Day(),
		slug,
	)

	f, err := os.Create(fn)

	if err != nil {
		return err
	}

	_, err = f.WriteString(
		fmt.Sprintf(
			"<!-- :metadata:\n\ntitle: %s\ntags: %s\n-- publishedAt: %s\nsummary:\n\n-->\n",
			title, tags, now.Format(time.RFC3339),
		),
	)

	if err != nil {
		return err
	}

	return sh.RunV("nvim", fn)
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
	}

	if err != nil {
		return err
	}

	lastInputModTime, err := lastModTime("./internal/view", articlesInPath)

	if err != nil {
		return err
	}

	if lastInputModTime == nil || outModTime == nil || lastInputModTime.After(*outModTime) {
		fmt.Println("View or CSS changes detected, running tailwind...")
		err := tailwindCmd("--postcss", "-i", tailwinInFile, "-o", tailwindOutFile, "--minify")

		if err != nil {
			return err
		}

		return sh.Run("touch", tailwindOutFile)
	}

	return nil
}

func convertArticles(reconvert bool) error {
	res, err := converter.Convert(articlesInPath, articlesOutPath, reconvert)

	if err != nil {
		return err
	}

	fmt.Print(res.String())
	return nil
}
