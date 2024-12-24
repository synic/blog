---
title: Doot, simple task runner for your projects
publishedAt: 2024-11-26T16:36:04-07:00
tags: [Programming, Python]
openGraph:
  image: "https://github.com/synic/doot/raw/aec35bbc68fc846c606ce04a14b9a1cce8c7ccdd/docs/images/thebestdoots.jpg"
summary: |
  Doot is a simple, zero dependency (except Python 3, which comes installed on
  most *nix operating systems) task runner. Similar to `make`, but meant to be
  used for non-C style projects.
---
Doot is a simple, zero dependency (except Python 3, which comes installed on
most *nix operating systems) task runner. Similar to `make`, but meant to be
used for non-C style projects.

# Getting Started

First,
[install](https://github.com/synic/doot/blob/main/README.md#installation)
doot. Then, in your project root directory, create a file (I usually call it
`do`, but it can be anything you want):

```python
#!/usr/bin/env python3

from doot import do

@do.task(allow_extra=True)
def bash(_, extra):
    """Bash shell on the web container."""
    do.run("docker exec -it api bash", extra)


@do.task()
def start():
    """Start all services."""
    do.run("docker-compose up -d")


@do.task()
def stop():
    """Stop all services."""
    do.run("docker-compose stop")


@do.task()
def dbshell():
    """Execute a database shell."""
    do.run("docker exec -it database psql -U myuser mydatabase")


@do.task()
def shell():
    """Open a django shell on the web container."""
    do.run("docker exec -it api django-admin shell")


@do.task(allow_extra=True)
def manage(_, extra):
    """Run a django management command."""
    do.run("docker exec -it api django-admin", extra)


@do.task(
    do.arg("-n", "--name", help="Container name", required=True),
    do.arg("-d", "--detach", help="Detach when running `up`", action="store_true"),
)
def reset_container(opt):
    """Reset a container."""
    do.run(f"docker-compose stop {opt.name}")
    do.run(f"docker-compose rm {opt.name}")

    extra = "-d" if opt.detach else ""
    do.run(f"docker-compose up {extra}")


if __name__ == "__main__":
    do.exec()
```

With this setup, you can run tasks like `./do manage`, `./do shell`, etc.

Running `./do -h` will show output like this:

```
Usage: ./do [task]

Available tasks:

  bash                   Bash shell on the web container
  dbshell                Execute a database shell
  manage                 Run a django management command
  reset-container        Reset a container
  shell                  Open a django shell on the web container
  start                  Start all services
  stop                   Stop all services
```

For more information, see the [GitHub
Repository](https://github.com/synic/doot)

## Acknowledgements

This project was named after our beloved Doots. She will be missed.

![Doots](https://github.com/synic/doot/raw/aec35bbc68fc846c606ce04a14b9a1cce8c7ccdd/docs/images/thebestdoots.jpg)
