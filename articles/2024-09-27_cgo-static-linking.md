<!-- :metadata:

title: Golang: Static linking with CGO and distroless
tags: Programming
published: 2024-09-27T11:10:00-0700
summary:

Do you use distroless? Have you tried to build a distroless docker image for
your Go project, only to see an error like `/bin/foo: no such file or
directory`? Maybe you spent a bunch of time trying to figure out why that file
isn't there, only to find out that it *IS* there, but you're still getting the
error?

This post is for you!
-->

One of the nice things about Go is that, by default, it compiles to a
statically linked binary with no external dependencies (including shared
libraries). This makes it super easy to deploy, and you can use a basically
empty docker image (like
[distroless](https://github.com/GoogleContainerTools/distroless)). However, if
you have `CGO_ENABLED=1` (which is required for some libraries, such as the C
based sqlite bindings), it will link everything statically except `libc` by
default. When using distroless, you usually compile on an image that has `libc`,
but then copy the binary to an image that doesn't, so while compilation doesn't
fail, the binary cannot find `libc` and will not run.

You can remedy this by statically linking `libc` by passing `-ldflags '-s -w
-linkmode external -extldflags "-static"'` to `go build` like so:

```bash
$ go build -a \
  -ldflags '-s -w -linkmode external -extldflags "-static"' \
  ./...
```

A couple things to note about this method:

* It will take quite a bit longer to compile
* The resulting file will be a lot larger
* This will NOT work if you use the Go plugin system (`dlopen` doesn't work)

Good luck!
