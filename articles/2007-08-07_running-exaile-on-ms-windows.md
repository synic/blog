<!-- :metadata:

title: Running Exaile on MS Windows
tags: Exaile
publishedAt: 2007-08-07T16:47:40-07:00
summary:

Exaile now runs on Windows, and runs quite well.  Thanks to the GStreamer and
Songbird people who recently ported GStreamer to Windows...

-->

<p>Exaile now runs on Windows, and runs quite well.  Thanks to the GStreamer
and Songbird people who recently ported GStreamer to Windows.</p>

<p>Functionally this is equivalent to the patched version in <a
href="http://www.exaile.org/trac/ticket/200">ticket #200</a>, only using
GStreamer instead of Windows Media Player.</p>

<p>As ticket #200 shows, it really doesn&#8217;t take much to make Exaile run
on Windows.  Most of it was already cross-platform, largely thanks to the
insistence on using <code>os.sep</code> and <code>os.path.join</code>.</p>

<b>How good</b>

<p>It&#8217;s as stable and fast as in Linux, and the majority of features work
perfectly.  Memory consumption is rather high, but then it&#8217;s high even in
Linux.</p>

<b>How bad</b>

<p>Encoding errors abound.  Don&#8217;t expect to be able to import most of
your library.</p>

<p>Anything to do with D-Bus has been chopped off and won&#8217;t work.  In
other words, Exaile doesn&#8217;t depend on D-Bus anymore; but that&#8217;s a
hack, not a feature.  I&#8217;m still hoping for D-Bus and dbus-python on
Windows so I can remove the hack.</p>

<p>Things go wrong when Exaile reloads open playlists from the previous
session.  From the symptoms, I&#8217;m guessing that Exaile thinks those tracks
are streams.</p>

<b>How it looks</b>

<strong>NOTE:</strong> _Unfortunately, this image has been lost to time._

<b>How</b>

<p>The changes I made are set for Exaile 0.2.11.  For now they&#8217;re only
available in trunk.</p>

<p>The requirements are basically identical to the ones mentioned <a
href="http://www.exaile.org/requirements">in Exaile&#8217;s Website</a>, except
that you don&#8217;t get to <code>apt-get</code> them.</p>

<ul>
<li><a href="http://www.python.org/">Python</a></li>
<li><a href="http://gladewin32.sourceforge.net/">
GTK+ and libglade for Windows</a></li>
<li><a href="http://www.pygtk.org/">PyGTK and its dependencies</a></li>
<li><a href="http://gstreamer.freedesktop.org/">
GStreamer, its plugins, and its Python bindings</a></li>
<li><a href="http://www.sacredchao.net/quodlibet/wiki/Development/Mutagen">
Mutagen</a></li>
<li>For Python < 2.5,
<a href="http://www.initd.org/tracker/pysqlite/wiki/pysqlite">pysqlite</a></li>
<li>For Python < 2.5,
<a href="http://effbot.org/zone/element-index.htm">ElementTree</a></li>
</ul>

<p>Right now binary GStreamer Python bindings are only available in the <span
class="caps">CVS</span> packages, and they&#8217;re compiled for Python 2.4 so
you&#8217;ll need the 2.4 version of everything.</p>

<div class='vimtip'>

<h3><strong>vim tip:</strong> <i>Wrapping text</i></h3>

<p>
Sometimes it&#8217;s useful to wrap text at a specified column.  You can use
Vim&#8217;s <code>gq</code> command and <code>textwidth</code>
(<code>tw</code>) option to do this.  Now, Vim has tons of options for wrapping
lines (it can even autowrap text while you&#8217;re typing) and <code>gq</code>
itself can be used in many different ways, but I normally use <code>gq</code>
to wrap a block of comments after I finished writing it.  Simply select the
lines you want to format (in Visual mode) and hit <code>gq</code>.  The
wrapping column is determined by <code>textwidth</code>, e.g. <code>:set
tw=80</code> to wrap at column 80.
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
