<!-- :metadata:

title: Supa Dupa Count
tags: Programming, Android, Scala
published: 2012-02-14T19:29:40-0700
summary:

[A few years
ago](http://synicworld.com/2009/12/22/supacount-published-android-market/)
I posted [SupaCount](http://bitbucket.org/synic/supacount) to the
Android market. Shortly after that, I stopped posting updates for three
years.

-->

[A few years
ago](http://synicworld.com/2009/12/22/supacount-published-android-market/)
I posted [SupaCount](http://bitbucket.org/synic/supacount) to the
Android market. Shortly after that, I stopped posting updates for three
years.

The reason for this? I lost my signing key. Even if I wanted to make
updates to it, I couldn't upload it to the same project on the Android
Market, at least, not without a lot of grief. So, I just avoided it...
though it wasn't a huge problem. Since I uploaded it, I haven't gotten
one error report in the market dashboard, so it must be working.
However, during that time using it myself, I noticed a few features I've
wanted here or there that would just make things easier.

So, over the weekend, I rewrote the entire app in Scala. Why a rewrite?
I dunno, just because I wanted an excuse to learn Scala and Scala for
Android. I've named this one "DupaCount". "SupaCount" will remain on the
market, but with a blurb that it's been replaced. During the rewrite, I
added the few features that I wanted:

-   If you reboot the phone while a timer is running, it will still
    alert as expected. This one really is more of a bug, but people
    probably didn't notice as most wouldn't make a timer that runs
    longer than a few hours. I've run timers for up to 2 months before,
    and though they always kept counting down, they wouldn't always
    alert if I had rebooted during that time.
-   There is an easy to see & access button at the bottom of the main
    screen called "Add Timer"
-   The "description" field in the Add Timer dialog has been made
    optional. If you leave it blank, it will default to a display of the
    timer time. This will make it much easier to add a timer.

Anyway, the app can be found here:
<https://market.android.com/details?id=com.synicworld.dupacount>, and the
sourcecode can be found here: <http://bitbucket.org/synic/dupacount>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
