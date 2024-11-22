
<!-- :metadata:

title: Handling multimedia keys in GNOME 2.18
tags: Exaile, Programming
publishedAt: 2007-05-16T12:31:00-0700
summary:

<p><span class="caps">GNOME</span> 2.18 introduced a new way for applications
to handle multimedia keys.  Previously you have to muck around with X events,
while now <span class="caps">GNOME</span> does it for you and you can get
control of mmkeys by requesting through D-Bus (to <span
class="caps">GNOME</span> Control Center&#8217;s Settings Daemon).  All
good until you realise that for cross-desktop support you still need the
old method anyway&#8212;unless, like Rhythmbox and Banshee, your app is
GNOME-based.</p>

<p>This article shows how we support both methods in Exaile, and how you can do
it, too.</p>

-->

<p><span class="caps">GNOME</span> 2.18 introduced a new way for applications
to handle multimedia keys.  Previously you have to muck around with X events,
while now <span class="caps">GNOME</span> does it for you and you can get
control of mmkeys by requesting through D-Bus (to <span
class="caps">GNOME</span> Control Center&#8217;s Settings Daemon).  All
good until you realise that for cross-desktop support you still need the
old method anyway&#8212;unless, like Rhythmbox and Banshee, your app is
GNOME-based.</p>

<p>Few people are aware of this change&#8212;most only realised after their
mmkeys didn&#8217;t work anymore (<a
href="https://bugs.launchpad.net/ubuntu/+source/rhythmbox/+bug/32917">Rhythmbox</a>,
<a href="http://bugzilla.gnome.org/show_bug.cgi?id=395433">Banshee</a>, <a
href="http://www.exaile.org/trac/ticket/399">Exaile</a>, <a
href="https://bugs.launchpad.net/ubuntu/+source/quodlibet/+bug/43464/comments/23">Quod
Libet</a>, <a href="http://bugs.kde.org/show_bug.cgi?id=145239">Amarok</a>).  I
couldn&#8217;t find any documentation regarding the new method, so I decided to
write one to help other media player writers.</p>

1. <strong>Bus name</strong>: org.gnome.SettingsDaemon
2. <strong>Object path</strong>: /org/gnome/SettingsDaemon
3. <strong>Interface</strong>: org.gnome.SettingsDaemon
4. <strong>Methods</strong>:
    * <strong>GrabMediaPlayerKeys</strong>(application, time)
        * `application` is simply the application name.
        * `time` determines which application name gets sent with the
          MediaPlayerKeyPressed signal if multiple applications have registered
          interest in controlling mmkeys.  Rhythmbox uses 0 for <em>time</em>
          (I think 0 simply pushes your application to a stack), so
          that&#8217;s probably what you should use as well.
    * <strong>ReleaseMediaPlayerKeys</strong>(application)
5. <strong>Signals</strong>:
    * <strong>MediaPlayerKeyPressed</strong>(application, key)
        * Must check that <em>application</em> is your application name.
        * `key` can be &#8216;Play&#8217;, &#8216;Pause&#8217;,
          &#8216;Stop&#8217;, &#8216;Next&#8217;, and &#8216;Previous&#8217;
          (at least in current gnomecc).

<p><a href="http://exaile.org/trac/changeset/2255">Exaile&#8217;s revision
2255</a> shows how we added common mmkeys support to Exaile.  I have made the
MmKeys class generic so other players can also use it if they want to.  Note
that I&#8217;ve added a &#8216;PlayPause&#8217; key to match the signal from
the mmkeys module.</p>

<p>This post has been in my blog queue for some time, and meanwhile in the Quod
Libet world, Joe Wreschnig has <a
href="https://bugs.launchpad.net/ubuntu/+source/quodlibet/+bug/43464/comments/21">refused</a>
to add <span class="caps">GNOME</span> mmkeys support to the core, Facundo
Batista <a
href="https://bugs.launchpad.net/ubuntu/+source/quodlibet/+bug/43464/comments/23">created
a patch</a> that does it anyway, and Ronny Haryanto <a
href="https://bugs.launchpad.net/ubuntu/+source/quodlibet/+bug/43464/comments/25">created
a plugin</a> based on the patch (and <a
href="http://ronny.haryan.to/archives/2007/05/11/d-bus-multimedia-keys-plugin-for-quod-libet/">blogged</a>
about it).  The plugin does a few things incorrectly, so <del>I&#8217;ll
notify</del> <ins>I&#8217;ve notified</ins> Ronny about it.</p>

<p>Acknowledgements:
<ul>
<li>Thanks to Albert Bicchi for alerting me in #exaile to this change.</li>
<li>The information here is based on the source codes of <a
href="http://svn.gnome.org/viewcvs/gnome-control-center/trunk/gnome-settings-daemon/gnome-settings-dbus.c?view=markup"><span
class="caps">GNOME</span> Control Center</a> and <a
href="http://svn.gnome.org/viewcvs/rhythmbox/trunk/plugins/mmkeys/rb-mmkeys-plugin.c?view=markup">Rhythmbox</a>,
and also a patch for <a
href="http://bugzilla.gnome.org/show_bug.cgi?id=395433#c1">Banshee</a>.</li>
</ul></p>

<p>No thanks to <span class="caps">GNOME</span> for being quiet about this mass
application breakage.</p>

<p><strong>Update 2007-05-16:</strong> Just found out that this has also been
fixed in <a href="http://www.listen-project.org/ticket/606">Listen</a> some
time ago, yay!</p>

<div class='vimtip'>

<h3><b>vim tip:</b> <i>Color schemes</i></h3>

<p>
Tired of looking at the same colours in Vim all the time?  Use the
<code>colorscheme</code> (<code>colo</code>) command.  For example, try
<code>:colorscheme evening</code> (if you want to revert, the default scheme is
named &#8220;default&#8221;).  The effects of colour schemes are more noticable
when you&#8217;re using GVim.

There are a few colour schemes included with Vim, normally located in
<code>/usr/share/vim/vimXX/colors/</code> (XX is 70 if you use Vim 7.0).
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
