<!-- :metadata:

title: The Purpose of Vim Tabs
tags: Programming, VIM
published: 2012-02-09T18:59:34-0700
summary:

Yesterday, after many years of using Vim, I've finally realized what the
purpose of Vim Tabs is.   My friend asked me to post this article, because she
was also stumped by their functionality.

-->

Yesterday, after many years of using Vim, I've finally realized what the
purpose of Vim Tabs is.   My friend asked me to post this article, because she
was also stumped by their functionality.

When they first came out, I tried them.  I didn't understand why they'd even
put tabs in, when there are far superior ways to navigate Vim buffers, but I
soon found out that tabs in Vim don't show buffers.  -They show, in essence,
separate instances of Vim.  Buffers open in one tab will not show up in another
tab-.  *EDIT*: As VimGuy points out, they do contain the same buffers.  Tabs
are more like separate layouts of your currently open buffers.

So what's the point?  With all the options to organize your buffers and manage
your files (BufExplorer, NERDTree, split windows, etc, etc), what are the tabs
useful for?

At my new job, there are 3 separate, large, codebases.  Yesterday I was working
in one system, coding away.  Probably 20 buffers open, all related to the same
project, when one of my coworkers asked me to help him in another one of our
systems.  I didn't want to lose my place, or vim "state", as it were, and
intuitively I created a tab, opened up a new instance of NERDTree in this tab,
and started opening new buffers related to the second system in this tab.  In
the first tab, all my open buffers and split windows were still there, waiting
for me to get back to them.

A minute later, it suddenly occurred to me, "So *THATS* what Vim tabs are for."

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
