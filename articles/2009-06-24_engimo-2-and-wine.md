<!-- :metadata:

title: Engimo 2 and Wine
tags: Linux
publishedAt: 2009-06-24T18:19:34-0700
summary:

After beating Enigmo in the iPhone, I was eager for more.  I'm a sucker for
these type of games.  I downloaded the demo of the Windows version of <a
href='http://www.pangeasoft.net/enigmo/downloads.html'>Enigmo 2</a>, and, to my
delight, it works great out of the box in <a
href='http://www.winehq.com'>Wine</a>!  It's got tons more doodads and etc than
the iPhone version, so I'd say it's a perfect upgrade path if you like that
game.

-->

After beating Enigmo in the iPhone, I was eager for more.  I'm a sucker for
these type of games.  I downloaded the demo of the Windows version of <a
href='http://www.pangeasoft.net/enigmo/downloads.html'>Enigmo 2</a>, and, to my
delight, it works great out of the box in <a
href='http://www.winehq.com'>Wine</a>!  It's got tons more doodads and etc than
the iPhone version, so I'd say it's a perfect upgrade path if you like that
game.

The installer works fine, but I recommend launching the Enigmo 2 application
after you have installed it with a command similar to this:

```bash
$ wine explorer /desktop=Enigmo,1024x768 \
  'c:\\Program Files\\Ideas From the Deep\\Enigmo 2 Supernova\\Enigmo 2.exe'
```

Enigmo 2 appears to have a maximum resolution of 1024x768, and this command
will allow you to run it in a Window that size.

It even works on my HP Mini 1030, which doesn't have the most fantastic video
card in the world :)

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
