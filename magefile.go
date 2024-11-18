//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	Default  = Dev
	cssTheme = "catppuccin-mocha"

	// paths
	binPath     = "bin"
	debugPath   = path.Join(binPath, "blog-debug")
	releasePath = path.Join(binPath, "blog-release")
	runCmd      = sh.RunCmd("go", "run")
	buildCmd    = sh.RunCmd("go", "build")
	tailwindCmd = sh.RunCmd(filepath.FromSlash("node_modules/.bin/tailwindcss"))
	minifyCmd   = sh.RunCmd(filepath.FromSlash("node_modules/.bin/css-minify"))

	// required command line tools (versions are specified in go.mod)
	tools = map[string]string{
		"air":   "github.com/air-verse/air",
		"templ": "github.com/a-h/templ/cmd/templ",
	}

	// aliases
	P = filepath.FromSlash
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

	if _, err := os.Stat(P("node_modules/.bin/tailwindcss")); err != nil {
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

	return buildCmd("-race", "-tags", "debug", "-o", debugPath, P("cmd/serve/serve.go"))
}

func (Build) Release() error {
	return buildCmd(
		"-tags",
		"release",
		"-ldflags",
		"\"-s -w\"",
		"-o",
		releasePath,
		P("cmd/serve/serve.go"),
	)
}

func Codegen() error {
	mg.Deps(Deps.Dev)

	err := sh.Run(path.Join(binPath, "templ"), "generate")

	if err != nil {
		return err
	}

	return tailwindCmd(
		"--postcss",
		"-i", P("internal/view/css/main.css"),
		"-o", P("cmd/serve/assets/css/main.css"),
		"--minify",
	)
}

type Articles mg.Namespace

func compileArticles(recompile bool) error {
	dir := P("./cmd/compile/")
	allFiles, err := os.ReadDir(dir)
	includeFiles := make([]string, 0, len(allFiles))

	if err != nil {
		return err
	}

	for _, f := range allFiles {
		if !strings.HasSuffix(f.Name(), "_test.go") {
			includeFiles = append(includeFiles, path.Join(dir, f.Name()))
		}
	}

	args := []string{"run"}
	args = append(args, includeFiles...)
	args = append(args, []string{
		"-i", "articles",
		"-o", P("cmd/serve/articles"),
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

	f, err := os.Create(P("internal/view/css/syntax.css"))

	if err != nil {
		return err
	}

	f.WriteString(data)
	f.Close()

	return minifyCmd("-f", P("internal/view/css/syntax.css"), "--output", P("cmd/serve/assets/css"))
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
