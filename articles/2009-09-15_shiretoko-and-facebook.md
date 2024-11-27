<!-- :metadata:

title: Shiretoko and Facebook
tags: Linux
publishedAt: 2009-09-15T15:17:29-07:00
summary:

For whatever reason (I can only assume it has to do with the product being
released /after/ Ubuntu 9.04 came out) the Ubuntu devs have decided to call
Firefox 3.5 "Shiretoko".  You can install alongside your regular Firefox
installation by typing @aptitude install firefox-3.5@ and you can run it by
typing @firefox-3.5@.  That's all well and good, except that the User-Agent
string is set to:

-->

For whatever reason (I can only assume it has to do with the product being
released /after/ Ubuntu 9.04 came out) the Ubuntu devs have decided to call
Firefox 3.5 "Shiretoko".  You can install alongside your regular Firefox
installation by typing @aptitude install firefox-3.5@ and you can run it by
typing @firefox-3.5@.  That's all well and good, except that the User-Agent
string is set to:

```
Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.1.2) Gecko/20090803 Ubuntu/9.04 (jaunty) Shiretoko/3.5.2
```

instead of:

```
Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.1.2) Gecko/20090803 Ubuntu/9.04 (jaunty) Firefox/3.5.2
```

Which breaks sites like Facebook (Facebook chat doesn't work).  To fix this, I
installed the "User Agent Switcher" extension, which, after a few days of use,
I have decided I don't like.  It doesn't remember user preferences, and you
have to reset your default User-Agent every time you start Firefox.

Looking around on the web, I found a much easier way to solve this problem.
Type `about:config` in your address bar, and search for the
`general.useragent.extra.firefox` entry, right click on it, choose "modify" and
change "Shiretoko" to "Firefox".

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
