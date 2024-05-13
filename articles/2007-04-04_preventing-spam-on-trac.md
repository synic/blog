<!-- :metadata:

title: Preventing spam on trac
tags: Miscellaneous, Linux
published: 2007-04-04T22:25:00-0700
summary:

Preventing spam on Trac...

-->

<b>Note:</b>  There is a <a
href='http://www.exaile.org/wiki/index.php?title=Misc:HackingTracAccountManager'>HOW-TO</a>
written up to show users how to deploy the information in this article.<br />

The <a href='http://www.exaile.org'>Exaile</a> website and project uses trac to
manage the homepage wiki, and ticket system.  Trac is a great product, but a
lot of users end up getting spam at some point.

We solved this problem for a while by installing two spam filtering plugins,
and installing <a href='http://trac-hacks.org/wiki/AccountManagerPlugin'>Trac's
Account Manager Plugin</a>.  This requires a user to register for the site
before they can post tickets and/or make comments on existing tickets (the
default behavior for trac is to allow anyone at all post).

Spam bots and otherwise started getting around this anyway.

One of the <a href='http://www.jokosher.org/'>Jokosher</a> developers is a user
of Exaile, and he found out about our dilemma.  He told us how they managed to
stop spam.

They simply added a field to their registration form that says "What is the
name of our audio editor?".  The answer is so obvious.  Jokosher is the name.
It's in the domain name of the site.  It's in the logo at the top of the page.

They claim this stopped spam all together, even better than an image captcha.
The idea is that even if a bot can't get around the registration, spammers will
employ human spammers for special exceptions.  An image captcha doesn't prevent
spam from these people at all... they can see the answer right in front of
their eyes.  However, it's probably not cost effective to make them look around
for something a legit user would already know.

I made some changes to Trac's Account Manager.  I added a field that says "What
is Adam's IRC nickname".  My nickname is in several obvious places on the
website, not to mention, there's a two click way to get into the IRC channel
via CGI::IRC.  In future releases of Exaile, a keyword will be available in the
Exaile about dialog which will be required to register.

I've written up a <a
href='http://www.exaile.org/wiki/index.php?title=Misc:HackingTracAccountManager'>howto</a>
on modifying the plugin.

<b>Note:</b> I only just deployed this on our site.  I'll keep you updated on
how it works out.

<div class='vimtip'>
  <h3><b>vim tip:</b> <i>:TOhtml</i></h3>
  <p>
    Sometimes you want to post your code on a website somewhere.  Sometimes it's
    nice to have this code display with the syntax highlighting you see in vim.
    Just type @:TOhtml@, and a new window window will open up with the html for the
    code you were looking at.  Note, that all the coloring will be the same
    coloring you have for the current vim colorscheme your using (even the
    background color), so be sure to pick a colorscheme that matches your site
    before you do this.
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
