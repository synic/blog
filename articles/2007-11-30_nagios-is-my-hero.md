---
title: Nagios is my hero
slug: nagios-is-my-hero
publishedAt: 2007-11-30T23:34:07-07:00
tags: [Linux, Administration]
summary: |
  Our company is growing fast.  When I started in March, 2003, we were hosting
  our website with some company in Texas.  I assumed the position of "Systems
  Administrator" (I knew Linux fairly well at the time, but I had no formal
  training in administration) when we moved from the
  hosting company to our first dedicated server colocated at <a
  href='http://www.xmission.com'>Xmission</a>.  Since that time, we moved from a
  shared rack to our own cabinet, now housing 10 servers (one is currently
  powered off, because if we turn it on, we risk blowing our power circuit...
  yes, we are waiting on a power upgrade)...
---
Our company is growing fast.  When I started in March, 2003, we were hosting
our website with some company in Texas.  I assumed the position of "Systems
Administrator" (I knew Linux fairly well at the time, but I had no formal
training in administration) when we moved from the
 hosting company to our
first dedicated server colocated at <a
href='http://www.xmission.com'>Xmission</a>.  Since that time, we moved from
a
 shared rack to our own cabinet, now housing 10 servers (one is currently
powered off, because if we turn it on, we risk blowing our power circuit...
yes, we are waiting on a power upgrade).<br><br>

As the number of servers we're using adds up, so does the stress of having to
manage it all.  There are a lot of little things to keep track of to make
sure
 everything is running smoothly, and doing so can sometimes be a lot of
work,
 especially when my official title here is "Programmer", and not
"Systems
 Administrator".  Earlier this month, I toured another companies data
center
 and I learned of something that, as of
 yesterday, is going to make
my life a <b>LOT</b> easier: <a
href='http://www.nagios.org'>Nagios</a>.<br><br>

From their site: "Nagios is an Open Source host, service and network
monitoring program."  To me, that doesn't quite sum up the capabilities of
this awesome system.  Here's what we now use it for:

<ul>
<li>Monitoring the RAID setups on every server (we currently
use 3ware, LSI, and software RAID setups; each one is monitored separately).
If a drive/array goes bad, our admin staff gets an email, and I get a text message.</li>
<li>Monitoring disk usage on every server; you can set a warning and a
critical threshold - for example, if disk usage goes over 80% you get a
warning, if it goes over 90% you get a critical notice.</li>
<li>Monitoring system load (thresholds here are also completely
customizable)</li>
<li>Monitoring MySQL replication status</li>
<li>Each server is PINGed periodically to make sure it is still up.  If a
machine goes down, a notification is sent out</li>
<li>Each web server is monitored to make sure i's receiving HTTP
connections</li>
<li>Each server is monitored to make sure it's receiving SSH
connections</li>
<li>Monitoring MySQL status (number of connections, slow queries,
etc)</li>
<li>Each service on the mail server (POP, IMAP, SMTP, the mail queue) is
monitored</li>
</ul>

For every item you monitor, you can specify when, how, and how
 often you get
notified of events (if at all).  You can also tell the system
 you want to be
notified when something returns to normal. <br>

Nagios uses plugins to monitor different items.  It comes with a bunch to
monitor commonly used systems and services (like http, ftp, disk usage, etc).
There is also a community website where people post plugins that they have
written here: <a
href='http://www.nagiosexchange.org'>http://www.nagiosexchange.org</a>.  If
it's not included, and you can't find one on nagiosexchange, they are super
easy to write in almost any language.<br><br>

Not only that,
Nagios gives you a nifty web interface to see what's going on.  Here are some
screenshots:

<b>The Status Map</b><br><br>
<img src='/images/statusmap.png'><br><br>
... I mostly included this screenshot because it looks cool.  It looks tons
cooler when you have hundreds of machines being monitored - if one of the
machines has a problem, the green circle around it shows as red.<br><br>
<b>The Host Overview</b><br><br>
<img src='/images/allservers.png'><br><br>
As you can see in that last screenshot, we currenly have one critical notice.
That's a degraded array that, even though I had written a script to check for
such a situation, we wouldn't have known about without installing Nagios (for
whatever reason, the script I wrote failed to notify us).  That array is
currently rebuilding thanks to this sweet system.<br><br>
In conclusion, if you're feeling growing pains at your company when it comes
to monitoring your servers, I highly recommend you give this a try.  I'll warn
you: It's not
exactly hard to set up, but it is rather tedious.  In my opinion, it's worth
the time you'll spend.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
