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

	Default    = Dev
	bin        = "./bin/blog"
	binDebug   = fmt.Sprintf("%s-debug", bin)
	binRelease = fmt.Sprintf("%s-release", bin)
	cssTheme   = "catppuccin-mocha"

	runCmd      = sh.RunCmd("go", "run")
	buildCmd    = sh.RunCmd("go", "build")
	tailwindCmd = sh.RunCmd("./node_modules/.bin/tailwindcss")
	minifyCmd   = sh.RunCmd("./node_modules/.bin/css-minify")

	goCommands = map[string]GoCommand{
		"air":   {url: "github.com/cosmtrek/air", version: airVersion},
		"templ": {url: "github.com/a-h/templ/cmd/templ", version: templVersion},
	}
)

func Dev() error {
	mg.Deps(Deps.Dev)

	return runGoCmdWith(map[string]string{"DEBUG": "true"}, "air")
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
	mg.Deps(Articles.Compile)

	return buildCmd("-race", "-tags", "debug", "-o", binDebug, "./cmd/serve/serve.go")
}

func (Build) Release() error {
	return buildCmd("-tags", "release", "-o", binRelease, "./cmd/serve/serve.go")
}

func Codegen() error {
	mg.Deps(Deps.Dev)

	err := runGoCmd("templ", "generate")

	if err != nil {
		return err
	}

	return tailwindCmd(
		"--postcss",
		"-i",
		"./internal/web/css/main.css",
		"-o",
		"./cmd/serve/assets/css/main.css",
		"--minify",
	)
}

type Articles mg.Namespace

func compileArticles(recompile bool) error {
	args := []string{"run", "cmd/compile/compile.go", "-i", "articles", "-o", "cmd/serve/articles"}

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

	f, err := os.Create("./internal/web/css/syntax.css")

	if err != nil {
		return err
	}

	f.WriteString(data)
	f.Close()

	return minifyCmd("-f", "./internal/web/css/syntax.css", "--output", "./cmd/serve/assets/css")
}

func Vet() error {
	return sh.Run("go", "vet", "./...")
}

func Test() error {
	return sh.Run("go", "test", "-race", "./...")
}

type GoCommand struct {
	url     string
	version string
}

// runGoCmd/runGoCmdWith  tries to run a go command
//
// It first tries to determine if the exutable is available in the path. If it
// is, it tries to run it. If it's not, it tries to install it and then run it.
func runGoCmdWith(env map[string]string, name string, args ...string) error {
	_, err := exec.LookPath(name)

	// if it is installed, try to run it
	if err == nil {
		err = sh.RunWithV(env, name, args...)

		if err != nil {
			return err
		}

		return nil
	}

	// if it's not installed, try to install it and then run it
	info, ok := goCommands[name]

	if !ok {
		return fmt.Errorf("command not defined: %s", name)
	}

	fmt.Printf("command \"%s\" not found, attempting to install...\n", name)
	err = sh.RunV("go", "install", fmt.Sprintf("%s@%s", info.url, info.version))

	if err != nil {
		return err
	}

	err = sh.RunWithV(env, name, args...)

	if err != nil {
		return err
	}

	return nil
}

func runGoCmd(name string, args ...string) error {
	return runGoCmdWith(nil, name, args...)
}
