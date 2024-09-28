<!-- :metadata:

title: Exaile 0.3 Roundup for October 29, 2008
tags: Exaile
published: 2008-10-29T18:28:29-0700
summary:

This week I've been focusing on album art support.  It was actually Aren's idea
that, from the beginning, we should create pluggable interfaces for things like
dynamic playlists, album art, lyrics, devices, and etc.  This means that album
art support wouldn't be hard coded to one place like Amazon, but anyone could
add a plugin to support album art from anywhere.  Currently, there are three
album art plugins:  local search, which checks local directories for image
names like 'cover.jpg' and etc., Amazon, and Last.FM.  The idea is that each
plugin will be searched in order until something is found.  Also Aren's idea:
the order that each plugin is searched should be user definable...

-->

<p>This week I've been focusing on album art support.  It was actually Aren's
idea that, from the beginning, we should create pluggable interfaces for things
like dynamic playlists, album art, lyrics, devices, and etc.  This means that
album art support wouldn't be hard coded to one place like Amazon, but anyone
could add a plugin to support album art from anywhere.  Currently, there are
three album art plugins:  local search, which checks local directories for
image names like 'cover.jpg' and etc., Amazon, and Last.FM.  The idea is that
each plugin will be searched in order until something is found.  Also Aren's
idea: the order that each plugin is searched should be user definable.</p>

<p>So, this week's tarball includes the following:</p>
 <p>
 <ul>
 <li>Album
Art Manager (it's basically just like the old one)</li>
 <li>Preferences page
to allow you to drag and drop the album art search methods</li>
 </ul>
</p>
 <p><img src='/media/images/screenie.jpg'
alt='Screenshot of album art preferences' border='0' /></p>
 <p>
 Also, I've
uploaded the current Exaile 0.3 <a href='http://www.exaile.org/doc'>API
Documentation</a>.  Note that it's probably going to change.
 </p>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
