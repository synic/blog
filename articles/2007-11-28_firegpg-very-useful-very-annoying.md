<!-- :metadata:

title: FireGPG - Very Useful, Very Annoying
tags: Miscellaneous
published: 2007-11-28T17:06:07-0700
summary:

<a href='http://firegpg.tuxfamily.org/'>FireGPG</a> is a <a
href='http://www.mozilla.com'>Firefox</a> extension that allows you to
encrypt/decrypt/sign any textbox on any page.  More than that, it integrates
with Gmail so you can easily send and receive encrypted email!  I've been using
it for this feature alone as soon as I found out about it's existence.  Now
there really is no reason to use an IMAP client in my opinion...

-->

<a href='http://firegpg.tuxfamily.org/'>FireGPG</a> is a <a
href='http://www.mozilla.com'>Firefox</a> extension that allows you to
encrypt/decrypt/sign any textbox on any page.  More than that, it integrates
with Gmail so you can easily send and receive encrypted email!  I've been using
it for this feature alone as soon as I found out about it's existence.  Now
there really is no reason to use an IMAP client in my opinion.<br><br>
 So
what's so annoying about it?  Instead of relying on Firefox's own extension
update notification system, FireGPG checks for updates every time you restart
Firefox, and pops up an intrusive dialog asking if you want to upgrade.  That
<i>could</i> be acceptable, except that FireGPG has updates ALL THE TIME.
Worse than that, sometimes this dialog is displayed when you click on a link
that opens a new window.  Ok, you finally bite the bullet and choose to update,
causing you to have to restart Firefox.  You are now interrupted again with
another intrusive dialog outlining the changes in this update.  In the event
that Firefox crashes (as it often does) before you quit (does anyone ever quit
Firefox?), the fact that you've been shown the changes dialog is not saved, and
you get to see it again when you restart.
 <br><br>If you're like me, you're
fairly annoyed at this point, and you want to turn the auto updates off.  So,
you go to the options, and click the "Disable FireGPG's auto updates".  Now you
get a confirmation dialog asking if you're sure you want to disable the auto
updates, and that they recommend that you don't.  If you click "Ok", you're
prompted with <i>another</i> dialog informing you that FireGPG is still in beta
and might have security issues (they don't mention this on their website in any
obvious manner) and asking, again, if you're sure you want to disable auto
updates.  I almost didn't want to; who wants potential security
issues?<br><br>
 Ok, I realize that security issues with GPG and encryption
might be a big problem.  I realize that with a new application that is
supposedly still in beta, not having automatic updates could result in a flow
of duplicate bugs on an issue tracker.  However, I feel that the behavior of
this application is totally unacceptable, causing me to want to not use it at
all. <br><br>
 Here are my suggestions to the authors of FireGPG (and other
extension developers):  Include the option to auto update, but have it disabled
by default:  Firefox has an update system that seems to work very well.  In the
case of FireGPG where "potential security issues" can in fact be a huge
problem, use a non-intrusive method to notify the user that there is a security
update.  An icon in the statusbar, or a system notification (like Gmail
Notifier) would be MUCH less annoying.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
