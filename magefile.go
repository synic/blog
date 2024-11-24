//go:build mage

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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
	binPath            = "bin"
	debugPath          = path.Join(binPath, "blog-debug")
	releasePath        = path.Join(binPath, "blog-release")
	outputCssPath      = "./assets/css"
	inputCssFile       = "./internal/view/css/main.css"
	outputCssFile      = fmt.Sprintf("%s/main.css", outputCssPath)
	syntaxCssFile      = "./internal/view/css/syntax.css"
	inputArticlesPath  = "./articles/"
	outputArticlesPath = "./assets/articles"
	buildInfoPath      = fmt.Sprintf("%s/internal", packageName)

	// commands
	runCmd      = sh.RunCmd("go", "run")
	buildCmd    = sh.RunCmd("go", "build")
	tailwindCmd = sh.RunCmd("node_modules/.bin/tailwindcss")
	minifyCmd   = sh.RunCmd("node_modules/.bin/css-minify")

	// required command line tools (versions are specified in go.mod)
	tools = map[string]string{
		"air":   "github.com/air-verse/air",
		"templ": "github.com/a-h/templ/cmd/templ",
	}
)

func Dev() error {
	mg.Deps(Deps.Dev)

	return sh.RunV(path.Join(binPath, "air"))
}

type Deps mg.Namespace

func (Deps) Dev() error {
	gobin, err := filepath.Abs(binPath)

	if err != nil {
		return err
	}

	for name, location := range tools {
		_, err = os.Stat(path.Join(binPath, name))

		if err == nil {
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		fmt.Printf("installing tool %s ...\n", location)
		err = sh.RunWithV(map[string]string{"GOBIN": gobin}, "go", "install", location)

		if err != nil {
			return err
		}

	}

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

	err := sh.Run(path.Join(binPath, "templ"), "generate", "-lazy")

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

	return minifyCmd("-f", syntaxCssFile, "--output", outputCssPath)
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

func maybeRunTailwind() error {
	var (
		lastInputModTime *time.Time = nil
		outModTime       *time.Time = nil
	)

	outInfo, err := os.Stat(outputCssFile)

	if err == nil {
		modTime := outInfo.ModTime()
		outModTime = &modTime
	}

	err = filepath.WalkDir("./internal/view", func(path string, d fs.DirEntry, err error) error {
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

		if lastInputModTime == nil || modTime.After(*lastInputModTime) {
			lastInputModTime = &modTime
		}

		return nil
	})

	if err != nil {
		return err
	}

	if lastInputModTime == nil || outModTime == nil || lastInputModTime.After(*outModTime) {
		fmt.Println("View or CSS changes detected, running tailwind...")
		err := tailwindCmd("--postcss", "-i", inputCssFile, "-o", outputCssFile, "--minify")

		if err != nil {
			return err
		}

		return sh.Run("touch", outputCssFile)
	}

	return nil
}

func convertArticles(reconvert bool) error {
	return converter.Convert(inputArticlesPath, outputArticlesPath, reconvert)
}
