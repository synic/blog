<!-- :metadata:

title: Exaile.org gets hacked!
tags: Exaile, Security
published: 2007-08-21T03:12:19-0700
summary:

I&#8217;ll admit.  I didn&#8217;t think this really happened to sites as
small as exaile.org, where they are running Linux and not too many services,
but I guess that type of thinking is what leads to these types of things...

-->

I&#8217;ll admit.  I didn&#8217;t think this really happened to sites as
small as exaile.org, where they are running Linux and not too many services,
but I guess that type of thinking is what leads to these types of things...

<p>After it happened I began reading blog posts and other material regarding
rootkits and security.  I was surprised, and rather distraught to find out that
this type of thing happens more than you might think, and that in most cases,
the owner of the machine in question doesn&#8217;t even know that it&#8217;s
happened.  I don&#8217;t know about you, but that worries me to death.  If
that&#8217;s the case, you might be wondering how I found out?</p>

<p>I&#8217;ll admit.  I didn&#8217;t think this really happened to sites as
small as exaile.org, where they are running Linux and not too many services,
but I guess that type of thinking is what leads to these types of things.</p>

<p>After it happened I began reading blog posts and other material regarding
rootkits and security.  I was surprised, and rather distraught to find out that
this type of thing happens more than you might think, and that in most cases,
the owner of the machine in question doesn&#8217;t even know that it&#8217;s
happened.  I don&#8217;t know about you, but that worries me to death.  If
that&#8217;s the case, you might be wondering how I found out? </p>

<p>First off, ssh starts acting weird.  Instead of the usual:</p>

```
The authenticity of host 'host.tld (3.344.33.67)'; can't be established.
RSA key fingerprint is e7:f9:61:d4:d3:f7:0c:42:3a:dc:39:4D:89:2b:7b:f8.

Are you sure you want to continue connecting (yes/no)?
```


<p>that you get when you&#8217;re trying to connect to a host you&#8217;ve
never connected to before, I was getting some error message regarding
&#8220;public keys&#8221; that I&#8217;ve never seen before, and even weirder,
I was getting it when trying to connect to hosts I had already accepted the
public key for.  I did some manpage reading, and wondered if I had
inadvertently switched my preferred authentication type when I was playing
around with <a
href='http://www-personal.umich.edu/~mressl/webshell/'>WebShell</a> (who knows,
maybe webshell is the cause of all my problems).  The error message I was
getting made me thing I had somehow switched to publickey authentication.  So,
I tried passing the <code>-o
PreferredAuthentications=password,keyboard-interactive</code>.  I got an error
saying &#8220;-o PreferredAuthentications is an invalid option&#8221;.</p>

<p>Now, I&#8217;m no <span class="caps">SSH</span> expert, but I had used that
option before, and it mentions is specifically in the manpage.  By this point,
I think I&#8217;ve totally broken something.  Again, I know I&#8217;m not an
<span class="caps">SSH</span> expert, so I take the Windows approach.  I decide
that I&#8217;m going to just remove ~/.ssh and reinstall openssh-client.
apt-get appears to remove everything correctly, but upon reinstall, it claims
it can&#8217;t overwrite /usr/bin/ssh because the &#8220;Operation is not
permitted&#8221;.  What?  I&#8217;m root.  I should be able to remove whatever
I want.  I jump into #openssh on freenode.  I ask them if there is some sort of
built in block for security regarding this.  </p>

<p>After fiddling around and with some suggestions from a user in that channel,
I find out that, for whatever reason, /usr/bin/ssh has been giving the ext2
attributes of &#8220;+uia&#8221;.  These attributes are set using the
<code>chattr</code> command.  From the chattr manpage:</p>

<ul>
<li><b>+a</b>: A file with the ‘a’ attribute set can only be open in append
mode  for   writing.    Only   the   superuser   or   a   process   possessing
the     CAP_LINUX_IMMUTABLE capability can set or clear this attribute. </li>
<li><b>+i</b>: A file with the ‘i’ attribute cannot be modified: it cannot be
deleted or  renamed,  no  link  can  be created to this file and no data can be
written to the file.  Only the superuser or a  process  possessing  the
CAP_LINUX_IMMUTABLE capability can set or clear this attribute.</li>
<li><b>+u</b>: When  a  file  with  the ‘u’ attribute set is deleted, its
contents are saved.  This allows the user to ask for its undeletion.
</ul>

<p>Ok&#8230;. great.  I had seen the +i attribute before and knew about chattr
from my gentoo days.  +i is used on the ext3 journal file.  No big deal, I
figure this is something debian has done for extra security.  I unset the
attributes with <code>chattr -iua</code>.  I try to reinstall openssh-client
again.  Same thing happens, only with ssh-add.  I go to /usr/bin and run
<code>chattr -iua ssh-*</code>.  The reinstall goes well this time.  Ssh
however, is still acting weird.  I&#8217;m getting a little weirded out now.  I
type <code>which ssh</code>.  What does it return?  <b>/usr/local/bin/ssh</b>.
<span class="caps">WHAT</span> <span class="caps">THE</span> HELL?!!?  Now I
know I&#8217;ve been hacked, there&#8217;s no reason anything would be in
/usr/local/anything on this machine.  I check /usr/local/bin and notice
there&#8217;s a myriad of ssh-related commands.</p>

<p>I start scouring the logs.  I see things like this: </p>

```
Aug 20 07:42:54 sp su: Successful su for nobody by root
Aug 20 07:42:54 sp su: + ??? root:nobody
Aug 20 07:42:54 sp su: (pam_unix) session opened for user nobody by (uid=0)
Aug 20 07:42:54 sp su: (pam_unix) session closed for user nobody
```

<p>I&#8217;m freaking.  Why is root suing to nobody?  (I later find out that
this is the cronjob that runs daily for the find utility, and it&#8217;s
totally normal).  I check other files in /usr/bin with lsattr.  Most of them
are set +iua.  I keep looking around&#8230; using <code>ps aux</code> and
looking for processes out of the ordinary.  I look around in my web
directories&#8230; and in http://www.exaile.org/files, I find a directory
that&#8217;s just named &#8220;&nbsp;&nbsp;&nbsp;&nbsp;&#8221;.  That&#8217;s
right, four spaces.  Inside is a directory called &#8220;www.irs.gov&#8221;,
and inside that&#8230;. what do I find?  A couple of html files and a php
script designed to convince whoever that they are on the irs.gov website, and
that they are entitled to a sum of over $100 as a return, so long as they fill
out their credit card information.  A quick look at the php script shows that
the information is just logging to a file.  Nothing is in the file&#8230;
yet.</p>

<p>I keep looking around&#8230; at this point, I&#8217;m sued to root.
Suddenly, I see a line that says &#8220;Killed&#8221;, and I&#8217;m back at my
user prompt.  Someone had killed my bash session.  The <code>who</code> command
showed only me.  I su back to root.  Killed again, and then once more, and
I&#8217;m totally kicked off the server.  I try to log back in, and my password
doesn&#8217;t work.  You can imagine how I felt.  I wonder to myself &#8220;can
I log into my webhosting companies support panel and shut the server
down?&#8221;  I try sshing as root, and am able to get in with the root
password (I know, I know, I shouldn&#8217;t have root logins enabled).  I
quickly type &#8216;halt&#8217;.  </p>

<p>Worst feeling ever.  I&#8217;m not worried about exaile.org or any of my
other opensource websites.  I&#8217;m kind of a backup nazi&#8230; but
it&#8217;s the other sites I&#8217;m hosting for my friends and their companies
that I&#8217;m worried about.  They can&#8217;t access their websites, they
can&#8217;t get their mail.</p>

<p>How does one recover from this?  I found out that my hosting company, for
$70 (which I didn&#8217;t have at the time), will set me up an entirely fresh
install of Debian on a new drive with my old drive as a slave.  It took me a
couple of days to get the money in my account to be able to do so.  During this
time is when I&#8217;m doing my research.</p>

<p>I got my hands on a couple rootkit tarballs.  I was able to read the source
of them, and found out, that at least for <b>shv5</b>, modifying
<code>chattr</code> (by renaming it and replacing it with a script) would
render the kit totally useless.  Probably other things like renaming tar and
etc would stop a script kiddie.</p>

<p><span class="caps">HOWEVER</span>, I realize that these are just little
steps (read: hacks) in securing a website.  I know that it starts with having a
secure setup in the first place&#8230; but again, I&#8217;m not an expert on
this either.  So how does a n00b like me secure their piddly servers?  I really
don&#8217;t know.  What I do know is that I was being very lax in the first
place.  I was running proftpd, and I didn&#8217;t need to be.  I was sending
ftp passwords to people on irc via message, though in plain text.  My web
directory permissions were often just &#8220;777&#8221; to make things easier
when I couldn&#8217;t figure out why something wasn&#8217;t working.  I had
heard that the imap and pop3 server (courier) I was using was insecure, but I
really wasn&#8217;t doing anything about it, and I wasn&#8217;t keeping up on
my security updates.  I was running a private <span class="caps">IRC</span>
server that was accessible publicly.  I&#8217;m running <span
class="caps">PHP</span>... I&#8217;ve always heard that&#8217;s a bad idea.
Probably a bunch of other things.</p>

<p>I&#8217;ve got my server back up, and most of the above problems are
corrected&#8230; all but courier, which is still running, and <span
class="caps">PHP</span>, which I&#8217;m still using.  I&#8217;ve got courier
limited via iptables to certain IPs because I really don&#8217;t know what else
to do.  All my users&#8217; current email depends on courier.</p>

<p>I know that I&#8217;m lucky in this case.  If the attacker had been anything
more than a 13 year old running scripts, I&#8217;d probably be hosting some
crazy phishing site and a number of <span class="caps">IRC</span> bots right
now without knowing it.  I&#8217;ve still got the old drive accessible as a
slave, but I don&#8217;t know what to do with it (is there a way to find out
how the attacker got in?)</p>

<p>Anyway.  That&#8217;s my story for the day.  If any of you have any
suggestions, I&#8217;d be glad to hear them.</p>

<p>I&#8217;m tired.  No vim tip for today&#8230; other than this:  Scrolling is
much faster in gVIM.  Just a random note.</p>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
