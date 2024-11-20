//go:build mage

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	Default     = Dev
	packageName = "github.com/synic/adamthings.me"
	cssTheme    = "catppuccin-mocha"

	// paths
	binPath     = "bin"
	debugPath   = path.Join(binPath, "blog-debug")
	releasePath = path.Join(binPath, "blog-release")
	runCmd      = sh.RunCmd("go", "run")
	buildCmd    = sh.RunCmd("go", "build")
	tailwindCmd = sh.RunCmd("node_modules/.bin/tailwindcss")
	minifyCmd   = sh.RunCmd("node_modules/.bin/css-minify")
	inputCss    = "./internal/view/css/main.css"
	outputCss   = "./assets/css/main.css"

	// misc
	buildInfoPath = fmt.Sprintf("%s/internal", packageName)

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
	mg.Deps(Articles.Compile)

	return buildCmd("-tags", "debug",
		fmt.Sprintf("-ldflags=-X %s.DebugFlag=true", buildInfoPath),
		"-o", debugPath,
		".",
	)
}

func (Build) Release() error {
	return buildCmd(
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

func compileArticles(recompile bool) error {
	dir := "./cmd/compile/"
	allFiles, err := os.ReadDir(dir)
	includeFiles := make([]string, 0, len(allFiles))

	if err != nil {
		return err
	}

	for _, f := range allFiles {
		if !strings.HasSuffix(f.Name(), "_test.go") && strings.HasSuffix(f.Name(), ".go") {
			includeFiles = append(includeFiles, path.Join(dir, f.Name()))
		}
	}

	args := []string{"run"}
	args = append(args, includeFiles...)
	args = append(args, []string{
		"-i", "articles",
		"-o", "articles/json",
		"-d",
	}...)

	if recompile {
		args = append(args, "-recompile", "-v")
	}

	return sh.RunV("go", args...)
}

func (Articles) Compile() error {
	return compileArticles(false)
}

func (Articles) Recompile() error {
	return compileArticles(true)
}

func Pygmentize() error {
	data, err := sh.Output("pygmentize", "-S", cssTheme, "-f", "html", "-a", ".chroma")

	if err != nil {
		return err
	}

	f, err := os.Create("internal/view/css/syntax.css")

	if err != nil {
		return err
	}

	f.WriteString(data)
	f.Close()

	return minifyCmd("-f", "internal/view/css/syntax.css", "--output", "cmd/serve/assets/css")
}

func Vet() error {
	return sh.Run("go", "vet", "./...")
}

func Test() error {
	return sh.Run("go", "test", "-race", "./...")
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

func maybeRunTailwind() error {
	var (
		lastInputModTime *time.Time = nil
		outModTime       *time.Time = nil
	)

	outInfo, err := os.Stat(outputCss)

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
		err := tailwindCmd("--postcss", "-i", inputCss, "-o", outputCss, "--minify")

		if err != nil {
			return err
		}

		return sh.Run("touch", outputCss)
	}

	return nil
}
