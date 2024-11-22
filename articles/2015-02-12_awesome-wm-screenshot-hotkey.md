<!-- :metadata:

title: Awesome-wm screenshot hotkey
tags: Programming, Linux
publishedAt: 2015-02-12T17:18:03-0700
summary:

At [Ender Labs](http://enderlabs.com), pretty much everyone but me uses a Mac.
In the last year or so, I've started to realize that there must be some new Mac
app or function that, via a hotkey, allows you to select an arbitrary region of
your screen, create a screenshot, and then automatically uploads it to a
hosting service. I know this, without doing any research, because in
irc/slack/gtalk I've started receiving screenshots as responses to questions I
ask. These screenshots arrive fairly quickly after said question is asked.

-->

At [Ender Labs](http://enderlabs.com), pretty much everyone but me uses a Mac.
In the last year or so, I've started to realize that there must be some new Mac
app or function that, via a hotkey, allows you to select an arbitrary region of
your screen, create a screenshot, and then automatically uploads it to a
hosting service. I know this, without doing any research, because in
irc/slack/gtalk I've started receiving screenshots as responses to questions I
ask. These screenshots arrive fairly quickly after said question is asked.

I rarely take screenshots (unless I absolutely have to), because my process
involves opening gimp, going to file->create->screenshot, selecting the "select
region to grab" option, setting the delay to 2 seconds so I can get the gimp
window out of the way, selecting the area I want to screenshot, exporting it to
a jpeg, and then uploading it to my server. It's enough hoops to jump through
that I usually just try to describe what I'm seeing rather than taking a
screenshot.

I liked this instant screenshot idea, and went looking for something like it
for Linux. I found a few apps, but nothing were exactly what I wanted, so I
decided to write my own. It's just a short, dirty python script that uses
[ImageMagick's](http://www.imagemagick.org/) `import` command to screenshot
a region, then automatically uploads the resulting jpeg to my server, and
copies the url to X's primary clipboard with
[xclip](http://linux.die.net/man/1/xclip). The python code is below.

```python
#!/usr/bin/env python

import os
import os.path
import uuid
import urlparse
fn = "/tmp/%s.jpg" % str(uuid.uuid4())
sn = os.path.basename(fn)
os.system("import %s" % fn)
os.system("notify-send \"Uploading %s\" ..." % sn)
os.system("scp %s synicworld.com:~/media/screenshots" % fn)
url = "http://synicworld.com/media/screenshots/%s" % sn
os.system("echo -n \"%s\" | xclip -selection primary" % url)
os.system("notify-send 'Copied to clipboard.'")
```

I saved the file to `~/bin/screenshot`, made it executable, and then put the
following in [awesome-wm's](http://awesome.naquadah.org/) rc.lua (keybindings
section):

```lua
awful.key(
  { modkey, "Shift" },
  "s",
  function()
    awful.util.spawn("/home/synic/bin/screenshot", false)
  end,
)
```

Now, all I have to do is press Win+Shift+s, select a region of my screen, and
the url is automatically copied to my clipboard. Neat!

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
