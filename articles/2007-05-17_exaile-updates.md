<!-- :metadata:

title: Exaile Updates
tags: Exaile
publishedAt: 2007-05-17T23:17:40-0700
summary:

Exaile 0.2.10 is well on it&#8217;s way...

-->

<p>Exaile 0.2.10 is well on it&#8217;s way.  </p>

<p>All plugins have been moved out of the regular Exaile repository, and into
their own.  The script located <a
href='http://exaile.org/plugins/plugins.py?version=trunk'>here</a> reads
information directly from <a
href='http://exaile.org/trac/browser/plugins/trunk'>Exaile&#8217;s trac source
browser</a> to keep up-to-date plugin list (trac watches our svn repository),
even if we&#8217;ve just committed changes.  </p>

<p>This will come to play in the new plugin manager, which I&#8217;ve started
to work on.  When the user opens the plugin manager, they will see an available
plugins tab.  When they click on that tab, an updated plugin list will be
downloaded from the <a
href='http://exaile.org/plugins/plugins.py?version=trunk'>plugin list
script</a> running on Exaile.org,</p>

<p>Currently, the user can install plugins from here.  Eventually, they will
also be able to remove and upgrade plugins from here.</p>

<p>Also, the pluggable radio system is done, and the shoutcast plugin has been
written.  Everything works as it used to, but new plugins will be able to be
written to get radio stations from other stream directories.</p>

<div class='vimtip'>

<h3><b>vim tip:</b> <i>Remote File Editing</i></h3>

<p>
You can edit files on a remote server that&#8217;s running an ssh daemon by
typing something like the following:

<code>:e scp://synic@jbother.org//etc/apt/sources.list</code>

he <code>:e</code>, of course, means edit file.  The <code>scp://</code>
means use the scp protocol (you can experiment with different protocols here.
http:// and ftp:// work, but obviously some of these protocols will be read
only).  The <code>synic@jbother.org</code> is &#8220;username at host&#8221;.
The following slash is required (and is usually a <b>:</b> on the command line,
but in vim, it&#8217;s a slash for a reasons I don&#8217;t know).  Everything
after that is the path of the file you want to edit on the remote server.

When you press enter, vim will prompt for the user&#8217;s password.  When
you save the file with <code>:w</code>, you will be prompted for the password
again.

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
