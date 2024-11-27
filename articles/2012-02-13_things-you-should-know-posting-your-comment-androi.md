<!-- :metadata:

title: Things you should know before posting your comment on the Android market
tags: Android
publishedAt: 2012-02-13T18:29:30-07:00
summary:

Yeah, ok, this post is more of a rant that anything.   I realize that people
actually reading this article probably aren't the ones that post my most
unfavorite comments on the market, but I don't care.  I'm posting it anyway.

-->

Yeah, ok, this post is more of a rant that anything.   I realize that people
actually reading this article probably aren't the ones that post my most
unfavorite comments on the market, but I don't care.  I'm posting it anyway.

# First and Foremost

The android market comments section is *not* a good place to post bug reports.
There is no way for the author of the app to get in touch with you if they need
more information about the bug, and no real reliable way for them to even
respond to your post.  Please please PLEASE try and find their bug tracker, if
there is one.

With that in mind, it's probably a better idea to use that same bug tracker to
file feature requests.  Especially if they are complicated.  Sure, it's fine
for something short and obvious, like "I wish you could change the font size",
etc (though I'd prefer that no one ever use the market for any feature
requests).

Do your research.  I often see comments like the following:

## "The latest update of this app starts up upon booting your phone, uninstalling"

Does this mean you just noticed that it requested boot permissions while
installing it?  Or you actually noticed it was running immediately after boot?
The world will never know, because, as stated above, the developer cannot
contact you based on your comment.  Regardless, noticing that the app requires
boot permissions or is running immediately after startup probably doesn't mean
what you think it means.

In the case of SupaCount, it starts on boot to reset the alarms that were wiped
out during the reboot.  This means it starts up, registers any alarms with the
Android Alarm Manager, and then quits.  However, because of the way Android
memory management works, it'll probably still show up in your Task Manager app
(more on these evil things later), because the OS keeps it around for a while,
just in case it's needed again.  If the OS needs this memory for something else
before that happens, it'll discard it, and that's when it will stop showing up
in your running tasks.

I imagine *most* apps are doing something similar to this.  A lot of apps use
the AlarmManager to wake the app up periodically so that it can perform some
background task.  For instance, the OkCupid app used to wake up every 15
minutes or so to check to see if you had any new messages.  In order for the
app to schedule these events with the AlarmManager, it needs to be started at
boot, otherwise, those events will never get registered.

## "One star.  App is constantly running in the background hogging all my memory"

Likely the app is just waking up to check something periodically, like I
mentioned above.  Even if it is actually running in the background all the
time, your memory is not what you need to be worried about.  If the OS needs
the memory and a certain app is hogging too much, that app will be killed.  No
exceptions.  Giving an app a bad rating or uninstalling it all together just
because you noticed it in your task manager is dumb.

What you should be worried about is your battery.  If an app is doing something
constantly in the background, it WILL eat away at your battery.  You can find
out what apps eat up the most battery by going to Settings->About
Phone->Battery Usage.

For the reasons I just listed above, you should uninstall your task manager.
There are countless articles out there about why they are bad (do a search if
you don't believe me), and the Android developers have stated themselves that
they aren't needed.  Android memory management works differently than it does
on your computer.  An app in the foreground has priority over any app in the
background.  If the OS needs memory for the app in the foreground, it will kill
apps in the background to get it.  In order to improve response times, android
may keep an app around in memory, even after you think you've exited, just in
case you need to open it again.  Like I said, though, if the OS needs that
memory, the app will be removed.

So, if you see an app running in the task manager, it really has no meaning to
you.  There are several reasons it could be there, and you killing it is more
likely to cause problems than anything else.  If you're worried about your
battery, check your battery usage section, and if there's an app behaving
badly, uninstall it.

## "I will give this app 5 stars if you enable installing to the SD card"

Uhg.  It's not as easy as that.  There is a very specific set of criteria that
determine if an app is installable to the SD card.  You can see that list here:
http://developer.android.com/guide/appendix/install-location.html#ShouldNot

A lot of it has to do with what happens when your SD card is unmounted (which
happens when you attach your phone to your computer).

In the case of SupaCount, it can't be installed to the SD card because the
startup apps are run before the SD card is attached, so any apps that need to
be run at startup cannot be installed on the SD card.


Don't be an asshole.  Especially if the app is free.  The author does not owe
you anything, and you freaking out the in comments is just going to get you
ignored.


A lot of this stuff probably applies to the iPhone App Store as well.  Just
have some common sense.  The rating and comment system is there simply so you
can state what you thought of the app, and so that other people can tell if the
app is worth installing.  It's not a place to report bugs or ask for features.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
