---
title: "Go: Meet glap, clap-style argument parsing"
slug: glap-clap-for-go
tags: [Programming, Go]
publishedAt: 2026-04-10T12:00:00-07:00
openGraph:
  title: "Golang: Meet glap, clap-style argument parsing"
  image: /static/images/articles/2026-04-10_glap-clap-for-go/glap-gopher.webp
summary: |
  Rustaceans have been spoiled by Rust's
  [clap](https://github.com/clap-rs/clap) (or, at least, I think they have).
  The derive macros, the ergonomics, the fact that you can describe an entire
  CLI with a struct and be done with it — it's really nice. I've never actually
  used (never tried Rust), but every time I see an example, I think "I wish Go
  had this." Go's standard `flag` package works, and
  [urfave/cli](https://github.com/urfave/cli) is the one I reach for when I
  need more, but neither has that declarative feel. So I wrote
  [glap](https://github.com/synic/glap).

  <img src="/static/images/articles/2026-04-10_glap-clap-for-go/glap-gopher.webp" alt="Go gopher" width="115" height="128" />

  Define your CLI with a struct and some tags:

  ```go
  type CLI struct {
      Config  string `glap:"config,short=c,required,help=Path to config file"`
      Verbose bool   `glap:"verbose,short=v,help=Enable verbose output"`
      Port    int    `glap:"port,short=p,default=8080,help=Port to listen on"`
      Output  string `glap:"output,short=o,possible=json|text|yaml,default=text,help=Output format"`
  }
  ```

  Parse it, use it. Subcommands, env var fallback, validators, groups,
  conflicts/requires, colored help, shell completions — all there. It's still
  beta, but it's working well for me.
---
Rustaceans have been spoiled by Rust's
[clap](https://github.com/clap-rs/clap) (or, at least, I think they have). The
derive macros, the ergonomics, the fact that you can describe an entire CLI
with a struct and be done with it — it's really nice. I've never actually used
(never tried Rust), but every time I see an example, I think "I wish Go had
this." Go's standard `flag` package works, and
[urfave/cli](https://github.com/urfave/cli) is the one I reach for when I need
more, but neither has that declarative feel. So I wrote
[glap](https://github.com/synic/glap).

<img src="/static/images/articles/2026-04-10_glap-clap-for-go/glap-gopher.webp" alt="Go gopher" width="115" height="128" />

## The struct tag API

This is the part that felt the most like clap. You describe your CLI as a
struct, put some tags on the fields, and hand it to glap:

```go
package main

import (
    "fmt"
    "os"

    "github.com/synic/glap"
)

type CLI struct {
    Config  string `glap:"config,short=c,required,help=Path to config file"`
    Verbose bool   `glap:"verbose,short=v,help=Enable verbose output"`
    Port    int    `glap:"port,short=p,default=8080,help=Port to listen on"`
    Output  string `glap:"output,short=o,possible=json|text|yaml,default=text,help=Output format"`
}

func main() {
    var cli CLI
    app := glap.New(&cli).Name("myapp").Version("1.0.0").About("Example")
    if _, err := app.Parse(os.Args[1:]); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    fmt.Println(cli.Config, cli.Port, cli.Output)
}
```

That's the whole program. You get `--help`, `--version`, validation of
`possible` values, defaults, required fields, short flags — all from the
tags.

## The builder API

If you'd rather not use reflection (ya know, if you're not into the whole
brevity thing), or you want to build commands dynamically
at runtime, there's a builder API that does the same thing:

```go
app := glap.NewCommand("myapp").
    Version("1.0.0").
    About("Builder API example").
    Arg(glap.NewArg("config").Short('c').Required(true).Help("Path to config file")).
    Arg(glap.NewArg("verbose").Short('v').Help("Enable verbose output")).
    Arg(glap.NewArg("port").Short('p').Default("8080").Help("Port to listen on")).
    Arg(glap.NewArg("output").Short('o').Default("text").
        PossibleValues("json", "text", "yaml").Help("Output format"))

matches, err := app.Parse(os.Args[1:])
```

Same CLI, but without those pesky struct tags. Pick whichever one fits
your project.

## What else it does

A few things I won't demo here but are worth knowing about:

* Nested subcommands (think `git remote add`)
* Env var fallback, so `--port` can be filled from `MYAPP_PORT`
* Custom validators
* Argument groups, with `conflicts_with` / `requires` relationships
* Count and append actions (`-vvv` for verbosity, repeatable flags)
* Auto-generated help with colored output
* Shell completions for bash, zsh, and fish.

## Status

It's still pre-1.0 and the API might shift a bit before I tag a stable
release, but I've been using it in my own projects and it's been solid. The
code, examples, and full docs live at
[github.com/synic/glap](https://github.com/synic/glap).

Give it a try!
