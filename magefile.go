//go:build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/synic/magex"
)

var (
	// dev dependencies
	airVersion = "v1.49.0"

	Default    = Dev
	bin        = filepath.FromSlash("bin/blog")
	binDebug   = fmt.Sprintf("%s-debug", bin)
	binRelease = fmt.Sprintf("%s-release", bin)
	cssTheme   = "catppuccin-mocha"

	runCmd      = sh.RunCmd("go", "run")
	buildCmd    = sh.RunCmd("go", "build")
	tailwindCmd = sh.RunCmd(filepath.FromSlash("node_modules/.bin/tailwindcss"))
	minifyCmd   = sh.RunCmd(filepath.FromSlash("node_modules/.bin/css-minify"))

	// aliases
	P = filepath.FromSlash
)

func Dev() error {
	mg.Deps(Deps.Dev)

	p, err := magex.MaybeInstallToolV("air", "github.com/cosmtrek/air", airVersion)

	if err != nil {
		return err
	}

	return sh.RunWithV(map[string]string{"DEBUG": "true"}, p)
}

type Deps mg.Namespace

func (Deps) Dev() error {
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

	return buildCmd("-race", "-tags", "debug", "-o", binDebug, P("cmd/serve/serve.go"))
}

func (Build) Release() error {
	return buildCmd("-tags", "release", "-o", binRelease, P("cmd/serve/serve.go"))
}

func Codegen() error {
	mg.Deps(Deps.Dev)

	version, err := magex.ModuleVersion("github.com/a-h/templ")

	if err != nil {
		return err
	}

	p, err := magex.MaybeInstallToolV("templ", "github.com/a-h/templ/cmd/templ", version)

	if err != nil {
		return err
	}

	err = sh.Run(p, "generate")

	if err != nil {
		return err
	}

	return tailwindCmd(
		"--postcss",
		"-i", P("internal/web/css/main.css"),
		"-o", P("cmd/serve/assets/css/main.css"),
		"--minify",
	)
}

type Articles mg.Namespace

func compileArticles(recompile bool) error {
	args := []string{
		"run",
		P("cmd/compile/compile.go"),
		"-i", "articles",
		"-o", P("cmd/serve/articles"),
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
	data, err := sh.Output("pygmentize", "-S", cssTheme, "-f", "html", "-a", ".chroma")

	if err != nil {
		return err
	}

	f, err := os.Create(P("internal/web/css/syntax.css"))

	if err != nil {
		return err
	}

	f.WriteString(data)
	f.Close()

	return minifyCmd("-f", P("internal/web/css/syntax.css"), "--output", P("cmd/serve/assets/css"))
}

func Vet() error {
	return sh.Run("go", "vet", "./...")
}

func Test() error {
	return sh.Run("go", "test", "-race", "./...")
}
