<!-- :metadata:

title: DAAP Music Sharing and Exaile
tags: Exaile, Multimedia
published: 2007-06-08T19:36:00-0700
summary:

<p>I&#8217;m pleased to announce that Exaile gained a new developer last week
&#8211; Aren Olson.  Aren has been contributing here and there for a while now,
and as of last week with his new plugin, <a
href='http://www.exaile.org/trac/browser/plugins/trunk/daap-share.py'>daap-share.py</a>.
What is <a
href='http://en.wikipedia.org/wiki/Digital_Audio_Access_Protocol'><span
class="caps">DAAP</span> </a> you might ask?</p>

-->

<p>I&#8217;m pleased to announce that Exaile gained a new developer last week
&#8211; Aren Olson.  Aren has been contributing here and there for a while now,
and as of last week with his new plugin, <a
href='http://www.exaile.org/trac/browser/plugins/trunk/daap-share.py'>daap-share.py</a>.
What is <a
href='http://en.wikipedia.org/wiki/Digital_Audio_Access_Protocol'><span
class="caps">DAAP</span> </a> you might ask?</p>

<p>Imagine that while at work, you could load up Exaile and connect to your
collection running on your machine at home, search through it just like you
search through your local collection, drag some tracks over to your playlist
and listen away.</p>

<p>You can now, and here&#8217;s a little how-to on getting it set up under
Ubuntu (Feisty Fawn).</p>

<p>I&#8217;m pleased to announce that Exaile gained a new developer last week
&#8211; Aren Olson.  Aren has been contributing here and there for a while now,
and as of last week with his new plugin, <a
href='http://www.exaile.org/trac/browser/plugins/trunk/daap-share.py'>daap-share.py</a>,
I asked him if he wanted commit access.  What is <a
href='http://en.wikipedia.org/wiki/Digital_Audio_Access_Protocol'><span
class="caps">DAAP</span> </a> you might ask?</p>

<p>Imagine that while at home, you could load up Exaile and connect to your
collection running on your machine at work, search through it just like you
search through your local collection, drag some tracks over to your playlist
and listen away.</p>

<p>With Aren&#8217;s plugin you can, and here&#8217;s a little how-to on
getting it set up under Ubuntu (Feisty Fawn).</p>

<p><span class="caps">DAAP</span> is the protocol that Apple uses in iTunes to
allow different users of iTunes on a local network to share their music with
each other. </p>

<p>Aren has created a Feisty repository for Exaile&#8217;s svn (which also
contains <a href='http://www.snorp.net/log/tangerine'>Tangerine</a> [a <span
class="caps">DAAP</span> server]), which he keeps pretty up to date, so first
things first, you&#8217;re going to want to add the following lines to your
/etc/apt/sources.list:</p>

```
deb http://download.tuxfamily.org/syzygy42 feisty exaile-svn
deb http://download.tuxfamily.org/syzygy42 feisty reacocard
```

<p>and then run the following:</p>

```
wget http://download.tuxfamily.org/syzygy42/8434D43A.gpg
sudo apt-key add 8434D43A.gpg
sudo apt-get update
sudo apt-get install exaile python-daap tangerine python-avahi
```

<p>Start up Exaile by going to Applications->Sound & Video->Exaile.  If this
your first time running Exaile, you will be prompted to add some directories to
your collection.  Do this, and wait for it to complete scanning.</p>

<p>Open up a terminal, and type <code>tangerine-properties</code>.  Check the
&#8220;Enable music sharing&#8221; box, the &#8220;Find in&#8221; radio butt,
and choose Exaile from the dropdown.  This will start Tangerine,  which will
begin serving up music from Exaile&#8217;s library.</p>

<p>In Exaile, go to Tools->Plugins and select the &#8220;Available
Plugins&#8221; tab.  When the list finishes loading (which can take a minute),
select the checkbox next to the &#8220;Music Sharing&#8221; plugin, and click
the &#8220;Install/Upgrade&#8221; button.  Close the plugin manager.</p>

<p>You will now see a new tab on the left called &#8220;Network&#8221;.  In the
dropdown at the top, you should see your local tangerine share (which will be
an IP and a port number, something like 66.79.34.24:52106), and then
&#8220;Custom Location&#8221;.   If you are working in an office with people
who are running iTunes, you might also see their shares.  You can connect to
any of these, so long as you have the password, and browse them like they were
your own collection.</p>

<p>Now, if you write down that IP and port for your local share, when you get
home from your office, you can follow this guide again to get Tangerine and
Exaile up and running, choose &#8220;Custom Location&#8221; from the dropdown
in network, and type the IP and port of your work share.  <span
class="caps">AWESOME</span>.</p>

<p>Note, this assumes that you don&#8217;t have a firewall at work&#8230; you
may need to forward a port or something, but I&#8217;ll leave that part to you.
The hard part is done&#8230;. thanks Aren!</p>

<div class='vimtip'>

<h3><b>vim tip:</b> <i>Setting up Vi keybindings in bash</i></h3>

<p>
If you&#8217;re like me, you&#8217;re constantly ending up with sequences
like this in your firefox address bar or at the command line in bash:
<code>ciw</code>.  That&#8217;s right, I&#8217;m trying to change the word that
my cursor is under.  It&#8217;s quite annoying.  Well, you can actually do this
in bash.  Open up a terminal, and type <code>set -o vi</code>.  You now have Vi
keybindings in <span class="caps">BASH</span>.  Note, that unlike in Vi, you
are in insert mode by default.  To enter command mode, you have to type escape.
Most commands work here&#8230; but there are a lot of exceptions.  If the
command you&#8217;re typing is getting complex, you can open an actual vi
screen for the command you&#8217;re typing by hitting the v key.
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
