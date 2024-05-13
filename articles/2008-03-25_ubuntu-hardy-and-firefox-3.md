<!-- :metadata:

title: Ubuntu Hardy and Firefox 3
tags: Linux
published: 2008-03-25T17:57:26-0700
summary:

Yesterday I decided to upgrade from Gutsy Gibbon to Hardy Heron (beta 1).  I
usually upgrade early, but my main motivation for this was to be able to
implement the new multimedia keys interface working in Exaile.  Apparently, the
Gnome crew decided to change this particular part of the DBus interface making
it backwards incompatable.  A description of the Exaile bug can be found <a
href='https://bugs.launchpad.net/exaile/+bug/191428'>here</a>...

-->

Yesterday I decided to upgrade from Gutsy Gibbon to Hardy Heron (beta 1).  I
usually upgrade early, but my main motivation for this was to be able to
implement the new multimedia keys interface working in Exaile.  Apparently, the
Gnome crew decided to change this particular part of the DBus interface making
it backwards incompatable.  A description of the Exaile bug can be found <a
href='https://bugs.launchpad.net/exaile/+bug/191428'>here</a>.<br /><br />
The upgrade was fairly painless, but here are a few of the gotchas I did
experience.<br /><br />

<b>Hardy Issues:</b><br /><br />

The only problem I had was sound not working once I finally booted into my
upgraded system.  I found that, by default, the -386 kernel was set as default
instead of the -generic kernel.  The -386 kernel does not come with snd-*
modules for some reason.  I edited `/boot/grub/menu.lst` and changed `default  0`
to `default  2`, which made it so the -generic kernel would load by default.
This fixed my sound issues.<br /><br />

<b>Firefox 3 Issues:</b><br /><br />

Other than missing extensions that are not yet compatable with Firefox 3, I've
really only had one really annoying problem.  This problem existed in Gutsy
when trying to upgrade to FF3, and remained once the upgrade to Hardy was
complete (Hardy comes with FF3 by default).  The problem is as follows:  If you
have Firefox on workspace one, and you click on a link in a terminal (or any
other application) on a different workspace, the entire firefox window moves to
that workspace (even if you have FF set to open links in a new tab).<br /><br
/>
 Googling for answers to this problem revealed nothing, and people in
irc.mozilla.org/#firefox couldn't really understand what I meant.  Asking in
irc.freenode.net/#ubuntu-us-ut, Christer Edwards informed me that he had the
same problem and that in a fresh install of Hardy (as opposed to an upgrade)
the problem did not exist.  After fiddling around for a bit, I found that a
simple `rm -rf ~/.mozilla` fixed the problem.<br /><br />
 I will keep posting
as other issues may or may not arise with the upgrade.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
