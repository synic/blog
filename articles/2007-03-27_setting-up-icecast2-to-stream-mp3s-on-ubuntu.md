---
title: Setting up icecast2 to stream mp3s on Ubuntu
publishedAt: 2007-03-27T21:29:00-07:00
tags: [Linux, Music]
summary: |
  I do this every time I have to reinstall the OS on my work machine... and as
  you saw in the previous article, it's about that time again.  I always forget
  to save the configuration files for it....
---
I do this every time I have to reinstall the OS on my work machine... and as
you saw in the previous article, it's about that time again.  I always forget
to save the configuration files for it....

Let's start with this:

```bash
$ sudo apt-get install icecast2 liblame-dev build-essential \
  libxml2-dev libshout3-dev
```

We will be using ices to stream to the server.  The package available in Ubuntu
does not have mp3 support, so we're going to compile a custom one against
liblame-dev.  You can download ices from here:
<a href='http://downloads.us.xiph.org/releases/ices/ices-0.4.tar.gz'>
  http://downloads.us.xiph.org/releases/ices/ices-0.4.tar.gz
</a>.

Run the following once you've got it:

```bash
$ tar -jxvf ices-0.4.tar.bz2
$ cd ices-0.4
$ ./configure --with-lame --prefix=/usr
$ make
$ sudo make install
```

Now comes the configuration.  Edit `/etc/icecast2/icecast.xml`.  Find the section
about authentication and change the source password.  This is the password that
ices will be using to communicate with the server.  Edit /etc/default/icecast2
and change "ENABLED" to true.

Start the server like so:

```bash
$ /etc/init.d/icecast2 start
```

You're going to need to create a log directory for ices:

```bash
$ sudo mkdir /var/log/ices && chown youruser /var/log/ices
```

There is a file in the ices source directory - `conf/ices.conf.dist`.  Edit it.
You'll want to change the stream name, genre, description, the name of the file
that will contain your playlist, whether or not you want the playlist to be
random (use 0 for no and 1 for yes), and password (use the same password that
you used in the icecast2 configuration file).  Save the file somewhere as
ices.conf.

Now to create your playlist.  Use the full path when using find like so:

```bash
$ find /home/synic/music -name "*.mp3" > ~/playlist.txt
```

Then, you can run ices using that configuration file:

```bash
$ ices -c ices.conf
```

You should then be able to connect to your station with your favorite player by
using the address http://localhost:8000/ices.  You can, of course, change
`localhost` to your ip.

That's it!

<div class='vimtip'>
<h3><b>vim tip:</b> <i>Split Windows</i></h3>

<p>
You can have as many split windows as you want in vim.  To split horizontally,
type <b>:sp</b>, optionally with the file you want to edit.  To split
vertically, type <b>:vs</b>.  To navigate between the windows, you can use
<b>Ctrl+W</b> followed by one of the directional keys <b>h j k</b>, and
<b>l</b>
</p>
</div>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
