<!-- :metadata:

title: Ubuntu 9.04 Update Notifications
tags: Linux
published: 2009-04-29T14:56:41-0700
summary:

Since I upgraded to Jaunty, I had been seriously irritated with their decision
to change the way updates are handled.

-->

<p>Since I upgraded to Jaunty, I had been seriously irritated with their
decision to change the way updates are handled.</p>

<p>Instead of the familiar update notification icon in your notification area,
the update window just pops up when a security update is available, or once a
week for everything else.  I'm not running OS X or Windows.  I'll update my
computer when I feel like it! </p>

<p>Fortunately, I was looking at the Jaunty release notes for an explanation of
this idiocy, and instead of finding good reasoning, I found a way to change it
back to the old behavior.</p>

`gconftool -s --type bool /apps/update-notifier/auto_launch false`

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
