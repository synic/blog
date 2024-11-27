<!-- :metadata:

title: Exaile LastFM Proxy Plugin
tags: Exaile, Programming, Music
publishedAt: 2008-04-17T18:12:32-07:00
summary:

A few versions ago, Exaile had Last.FM streaming support using <a
href='http://svn.base-art.net/public/elisa/lastfm/lastfm_src.py'>LastFMSource</a>,
a pygst plugin created by Philippe Normand of <a
href='http://elisa.fluendo.com'>Elisa</a>. It worked... sort of.  Every other
time you tried to connect to LastFM station, Gstreamer would lock up entirely,
taking Exaile out with it.  Not being able to fix this problem, it was
eventually removed from Exaile entirely.<br><br>
 Enter: <a
href='http://vidar.gimp.org/?page_id=50'>LastFM Proxy</a>...

-->

A few versions ago, Exaile had Last.FM streaming support using <a
href='http://svn.base-art.net/public/elisa/lastfm/lastfm_src.py'>LastFMSource</a>,
a pygst plugin created by Philippe Normand of <a
href='http://elisa.fluendo.com'>Elisa</a>. It worked... sort of.  Every other
time you tried to connect to LastFM station, Gstreamer would lock up entirely,
taking Exaile out with it.  Not being able to fix this problem, it was
eventually removed from Exaile entirely.

Enter: <a
href='http://vidar.gimp.org/?page_id=50'>LastFM Proxy</a>.  This is a program
written in python created to connect to LastFM and start streaming the music to
a proxy that you can connect to on your local machine using any music player
that supports HTTP streaming.

After a bit of hacking (and really,
this is some seriously hackish stuff), I've created a plugin for Exaile called
"LastFM Radio" that (mostly) seemlessly integrates LastFMProxy into Exaile.  To
the user, it appears to just be a Radio Panel plugin like the current
"Shoutcast Radio" plugin.  The user just clicks on the station they want to
listen to, and it starts playing.  They can "Skip", "Ban", or "Love" tracks
just like in the LastFM native player.

It still needs some work, but overall, I'm pretty pleased with how well it
works.  Give it a try!

Note:  You must be using the latest bzr version of Exaile for this plugin to
work.  You can get instructions on doing that from Exaile's <a
href='http://www.exaile.org/downloads'>downloads</a> page.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
