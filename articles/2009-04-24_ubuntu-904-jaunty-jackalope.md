<!-- :metadata:

title: Ubuntu 9.04  - Jaunty Jackalope
tags: Linux
publishedAt: 2009-04-24T00:59:41-07:00
summary:

Following the tradition of retarded naming schemes, Jaunty Jackalope was released today...

-->

<p>Following the tradition of retarded naming schemes, Jaunty Jackalope was
released today.  I've successfully upgraded 4 machines (three from aptitude
safe-upgrades, and one from a fresh install).  Here are some notes I've
gathered so far:</p>
<p><ul>
<li>My home machine would not load X, complaining about the nvidia drivers.  In
order to fix this, I had to install nvidia-glx-180</li>
<li>My home machine is using a prism54 card, which did not work after the
upgrade.  After fiddling around with it for a while, I found that blacklisting
the p54pci module (which worked prior to the upgrade) and loading the prism54
module instead worked just fine</li>
<li>My EeePC 1000h (which was installed from scratch) worked perfectly once the
installation completed.  This hasn't ever happened before with previous
versions of Ubuntu.  I usually needed to manually install wireless and nic
drivers.</li>
<li>One of the machines I upgraded failed near the end with a problem
concerning dbus and hal.  An "apt-get install --reinstall hal" fixed the
problem</li>
<li>The application menu editor no longer works whatsoever.  Right click on
"Applications", go to edit menus, and nothing happens.  This is the case on all
of my machines.</li>
</ul>
</p>
<p>All in all, the upgrade process hasn't been that bad.  There are some cool
new graphics (the new usplash and gdm themes look much better).  There haven't
been any "holy crap, look at that" features that I've seen yet, but it's only
been less than a day.</p>
<p>On another note, if you're in the software development business and you use
Ubuntu on your production servers, don't plan any of your own releases near any
of Ubuntu releases.  Trying to install some needed packages on our servers
proved difficult with everyone else hoarding the bandwidth trying to download
the new version.</p>

<h3>Update: Sat Apr 25, 9:39AM</h3>

<p>I've upgraded my last machine from Intrepid to Jaunty, my EeeBox.  Flawless
victory, upgrade went without a hitch.  I'd say that all in all, this has been
the easiest upgrade for me. </p>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
