---
title: "Bazitis:  Gitosis for Bzr (Baazar)"
publishedAt: 2008-11-26T00:24:01-07:00
tags: [Programming, Python, Administration]
summary: |
  At work, I use git.  For Exaile I use bzr.  I like them both quite a bit.  At
  work, we use <a
  href='http://eagain.net/gitweb/?p=gitosis.git;a=summary'>Gitosis</a> to manage
  our repositories and I have to say, it's pretty damn cool.  Nothing quite like
  this exists for bzr, so I ported Gitosis to bzr and called it Bazitis.  The
  launchpad project page can be found <a
  href='http://www.launchpad.net/bazitis'>here</a>.  Here are the instructions on
  how to use Bazitis:
---
<p>At work, I use git.  For Exaile I use bzr.  I like them both quite a bit.
At work, we use <a
href='http://eagain.net/gitweb/?p=gitosis.git;a=summary'>Gitosis</a> to manage
our repositories and I have to say, it's pretty damn cool.  Nothing quite like
this exists for bzr, so I ported Gitosis to bzr and called it Bazitis.  The
launchpad project page can be found <a
href='http://www.launchpad.net/bazitis'>here</a>.  Here are the instructions on
how to use Bazitis:</p>

<p>First off, I'd like to give credit to some people.  Tommi Virtanen is the
author of Gitosis.  Bazitis is a copy of the Gitosis code, all except for the
parts where I had to get a little hacky with bzrlib.  His website is <a
href='http://eagain.net'>http://eagain.net</a>.  The other person I'd like to
thank is Garry Dolley, who wrote a great blog post on how to use Gitosis, which
can be found here: <a
href='http://scie.nti.st/2007/11/14/hosting-git-repositories-the-easy-and-secure-way'>hosting
git repositories the easy and secure way</a>.  Garry has given me permission to
copy his instructions and modify them for Bazitis, as long as I give him kudos,
which I have done in this paragraph.  Thanks guys!</p>

<p><h2>Install bazitis</h2></p>

<p>bazitis is a tool for hosting bzr repositories (I'm repeating myself for
those who skim :)</p>

<p>The first thing to do is grab a copy of bazitis and install it on your
server:</p>

```bash
$ cd ~/src
$ bzr branch lp:bazitis
```

<p>Next, install it like so:</p>

```bash
$ cd bazitis
$ python setup.py install
```

<p>Don't use --prefix unless you like self-inflicted pain. It is possible to
install bazitis in a non-standard location, but it's not nice. Read the Caveats
section at the bottom and then come back here.</p>

<p>If you get this error:</p>

```bash
-bash: python: command not found
```

<p>or</p>

```
Traceback (most recent call last):
  File "setup.py", line 2, in ?
    from setuptools import setup, find_packages
ImportError: No module named setuptools
```

<p>You have to install Python setuptools. On Debian/Ubuntu systems, it's just:</p>

```bash
$ sudo apt-get install python-setuptools
```

<p>For other systems, someone tell me or leave a comment, so I can update this
section and improve this tutorial.</p>

<p>The next thing to do is to create a user that will own the repositories you
want to manage. This user is usually called bzr, but any name will work, and
you can have more than one per system if you really want to. The user does not
need a password, but does need a valid shell (otherwise, SSH will refuse to
work).</p>

```bash
$ sudo adduser \
    --system \
    --shell /bin/sh \
    --gecos 'bzr version control' \
    --group \
    --disabled-password \
    --home /home/bzr \
    bzr
```

<p>You may change the home path to suit your taste. A successful user creation
will look similar to:</p>

```
Adding system user `bzr'...
Adding new group `bzr' (211).
Adding new user `bzr' (211) with group `bzr'.
Creating home directory `/home/bzr'.
```

<p>You will need a public SSH key to continue. If you don't have one, you may
generate one on your local computer:</p>

```bash
$ ssh-keygen -t rsa
```

<p>The public key will be in $HOME/.ssh/id_rsa.pub. Copy this file to your
server (the one running bazitis).</p>

<p>Next we will run a command that will sprinkle some magic into the home
directory of the bzr user and put your public SSH key into the list of
authorized keys.</p>

```bash
$ sudo -H -u bzr bazitis-init < /tmp/id_rsa.pub
```

<p>id_rsa.pub above is your public SSH key that you copied to the server. I
recommend you put it in /tmp so the bzr user won't have permission problems
when trying to read it. </p>

<p>Here some cool magic happens. Run this on your local machine:</p>

```bash
$ bzr branch bzr+ssh://bzr@YOUR_SERVER_HOSTNAME/bazitis-admin
$ cd bazitis-admin
```

<p>You will now have a bazitis.conf file and keydir/ directory:</p>

```bash
~/dev/bazitis-admin (master) $ ls -l
total 8
-rw-r--r--   1 garry  garry  104 Nov 13 05:43 bazitis.conf
drwxr-xr-x   3 garry  garry  102 Nov 13 05:43 keydir/
```

<p>This repository that you just cloned contains all the files (right now, only
2) needed to create repositories for your projects, add new users, and defined
access rights. Edit the settings as you wish, commit, and push. Once pushed,
bazitis will immediately make your changes take effect on the server. So we're
using Bzr to host the configuration file and keys that in turn define how our
Bzr hosting behaves. That's just plain cool.</p>

<p>From this point on, you don't need to be on your server. All configuration
takes place locally and you push the changes to your server when you're ready
for them to take effect.</p>

<p><h2>Creating new repositories</h2></p>

<p>This is where the fun starts. Let's create a new repository to hold our
project codenamed FreeMonkey.</p>

<p>Open up bazitis.conf and notice the default configuration:</p>

```
[bazitis]

[group bazitis-admin]
writable = bazitis-admin
members = jdoe
```

<p>Your "members" line will hold your key filename (without the .pub extension)
that is in keydir/. In my example, it is "jdoe", but for you it'll probably be
a combination of your username and hostname.</p>

<p>To create a new repo, we just authorize writing to it and push. To do so,
add this to bazitis.conf:</p>

```
[group myteam]
members = jdoe
writable = free_monkey
```

<p>This defines a new group called "myteam", which is an arbitrary string.
"jdoe" is a member of myteam and will have write access to the "free_monkey"
repo.</p>

<p>Save this addition to bazitis.conf, commit and push it:</p>

```
$ bzr commit -m "Allow jdoe write access to free_monkey"
$ bzr push bzr+ssh://bzr@YOUR_SERVER_HOSTNAME/bazitis-admin
```

<p><b>Note:</b> You only have to add the path to the bazitis-admin repo the
first time you push.  After that, it will be remembered and you can just type
"bzr push"</p>

<p>Now the user "jdoe" has access to write to the repo named "free_monkey", but
we still haven't created a repo yet. What we will do is create a new repo
locally, and then push it:</p>

```bash
$ mkdir free_monkey
$ cd free_monkey
$ bzr init

# do some work, bzr add and commit files

$ bzr push bzr+ssh://bzr@YOUR_SERVER_HOSTNAME/free_monkey
```

<p>With the final push, you're off to the races. The repository "free_monkey"
has been created on the server (in /home/bzr/repositories) and you're ready to
start using it like any ol' bzr repo.</p>

<p><h2>Adding users</h2></p>

<p>The next natural thing to do is to grant some lucky few commit access to the
FreeMonkey project. This is a simple two step process.</p>

<p>First, gather their public SSH keys, which I'll call "alice.pub" and
"bob.pub", and drop them into keydir/ of your local bazitis-admin repository.
Second, edit bazitis.conf and add them to the "members" list.</p>

```bash
$ cd bazitis-admin
$ cp ~/alice.pub keydir/
$ cp ~/bob.pub keydir/
$ bzr add keydir/alice.pub keydir/bob.pub
```

<o>Note that the key filename must have a ".pub" extension.</p>

<p>bazitis.conf changes:</p>

```diff
 [group myteam]
- members = jdoe
+ members = jdoe alice bob
  writable = free_monkey
```

<p>Commit and push:</p>

```bash
$ bzr commit -m "Granted Alice and Bob commit rights to FreeMonkey"
$ bzr push
```

<p>That's it. Alice and Bob can now clone the free_monkey repository like so:</p>

```bash
$ bzr branch bzr+ssh://bzr@YOUR_SERVER_HOSTNAME/free_monkey
```

<p>Alice and Bob will also have commit rights. </p>

<p><h2>Limitations</h2></p>

<p><ul>
<li>Currently, bazitis doesn't support everything that gitosis does, like
public readonly access.  This is planned for the future.  </li>
 <li>I haven't
tested bazitis with shared bzr repositories.  I have no idea if it will work.
If you try this, let me know how it goes.</li>
 <li>Bazitis works best with
bzr 1.9.  It works with earlier versions, but if you try to access a repository
that you do no have permission for, a huge ugly exception is thrown that would
probably lead a user to think something is wrong with bzr.  This is handled a
lot better in later versions of bzr.</li>
</ul>
</p>

<p>And that's all.  Let me know how it works out for you!</p>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
