//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	// dev dependencies
	airVersion   = "v1.49.0"
	templVersion = "v0.2.778"

	bin        = "./bin/blog"
	binDebug   = fmt.Sprintf("%s-debug", bin)
	binRelease = fmt.Sprintf("%s-release", bin)
	cssTheme   = "catppuccin-mocha"
	Default    = Dev
)

func Dev() error {
	mg.Deps(Deps.Dev)

	_, err := exec.LookPath("air")

	if err == nil {
		err = sh.RunV("air")
	} else {
		err = sh.RunV("go", "run", fmt.Sprintf("github.com/cosmtrek/air@%s", airVersion))
	}
	return err
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

	return sh.Run(
		"go", "build",
		"-race",
		"-tags", "debug",
		"-o", binDebug,
		"./cmd/serve/serve.go",
	)
}

func (Build) Release() error {
	return sh.Run(
		"go", "build",
		"-tags", "release",
		"-o", binRelease,
		"./cmd/serve/serve.go",
	)
}

func Codegen() error {
	mg.Deps(Deps.Dev)

	if err := sh.Run("templ", "generate"); err != nil {
		return err
	}

	return sh.RunV(
		"./node_modules/.bin/tailwindcss",
		"--postcss",
		"-i", "./internal/web/css/main.css",
		"-o", "./cmd/serve/assets/css/main.css",
		"--minify",
	)
}

type Articles mg.Namespace

func compileArticles(recompile bool) error {
	args := []string{
		"run",
		"cmd/compile/compile.go",
		"-i", "articles",
		"-o", "cmd/serve/articles",
	}

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
	data, err := sh.Output(
		"pygmentize",
		"-S", cssTheme,
		"-f", "html",
		"-a", ".chroma",
	)

	if err != nil {
		return err
	}

	f, err := os.Create("./internal/web/css/syntax.css")

	if err != nil {
		return err
	}

	f.WriteString(data)
	f.Close()

	return sh.RunV(
		"./node_modules/.bin/css-minify",
		"-f", "./internal/web/css/syntax.css",
		"--output", "./cmd/serve/assets/css",
	)
}

func Vet() error {
	return sh.Run("go", "vet", "./...")
}

func Test() error {
	return sh.Run("go", "test", "-race", "./...")
}
