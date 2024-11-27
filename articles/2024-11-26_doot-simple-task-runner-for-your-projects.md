<!-- :metadata:

title: Doot, simple task runner for your projects
tags: Programming, Python
publishedAt: 2024-11-26T16:36:04-07:00
ogImage: https://github.com/synic/doot/raw/aec35bbc68fc846c606ce04a14b9a1cce8c7ccdd/docs/images/thebestdoots.jpg
summary:

Doot is a simple, zero dependency (except Python 3, which comes installed on
most *nix operating systems) task runner. Similar to `make`, but meant to be
used for non-C style projects. Comes out of the box with simple docker support.

-->

Doot is a simple, zero dependency (except Python 3, which comes installed on
most *nix operating systems) task runner. Similar to `make`, but meant to be
used for non-C style projects. Comes out of the box with simple docker support.

# Installation

I prefer using the [Zero
Install](https://github.com/synic/doot/tree/aec35bbc68fc846c606ce04a14b9a1cce8c7ccdd?tab=readme-ov-file#zero-install-option)
option, as doing it this way means that your coworkers don't have to install
anything extra to get their runner working (assuming they already have Python
installed, which is usually true).

Alternatively, you can install it as a library:

```bash
$ pip install git+https://github.com/synic/doot
```

# Getting Started

In your project root directory, create a file (I usually call it `do`, but it
can be anything you want):

```python
#!/usr/bin/env python

import os

import doot as do

@do.task(passthrough=True)
def bash(opts):
    """Bash shell on the web container."""
    do.crun("bash", opts.args)


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
    do.crun("psql -U myuser mydatabase", container="database")


@do.task()
def shell():
    """Open a django shell on the web container."""
    do.crun("django-admin shell")


@do.task(passthrough=True)
def manage(opts):
    """Run a django management command."""
    do.crun("django-admin", opts.args)


@do.task(
    do.opt("-n", "--name", help="Container name", required=True),
    do.opt("-d", "--detach", help="Detach when running `up`", action="store_true"),
)
def reset_container(opts):
    """Reset a container."""
    do.run(f"docker-compose stop {opts.name}")
    do.run(f"docker-compose rm {opts.name}")

    extra = "-d" if opts.detach else ""
    do.run(f"docker-compose up {extra}")


if __name__ == "__main__":
    do.main(default_container="web")
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
