---
title: My Introduction
publishedAt: 2007-04-04T18:48:00-07:00
tags: [Random]
summary: |
  Adam has kindly let me post here...
---
<p>Adam has kindly let me post here, and it&#8217;s only fitting that I use my
first post to introduce myself. I&#8217;m Johannes Sasongko, I live in
Brisbane, Australia, and I&#8217;m currently studying for my IT degree.
I&#8217;m involved in a few open-source projects at the moment: Ruby.<span
class="caps">NET</span> (as tester), Exaile (as programmer), and K-Meleon (as
bug triager / <span class="caps">BTS</span> cleaner, though the quality of the
recent versions has left me with little work).</p>

<p>Following Adam, and in the spirit of this site&#8217;s domain name, I will
of course post Vim tips after my posts. And since this is his blog, I&#8217;ll
be writing about things that he would be interested in. I have little interest
in leaving &#8216;Web-space junk&#8217; anyways; the signal-to-noise ratio of
the &#8216;Net is already bad as it is.</p>

<div class="vimtip">

<h3>
<b>vim tip:</b> <i>Indentation</i>
</h3>

<p>
Vim has a few settings you can use to configure the indentation of your
code. <code>expandtab</code> (<code>et</code>) is a boolean setting that does
what it says, it specifies whether to expand tabs into spaces; set
<code>noexpandtab</code> (<code>noet</code>) to disable it.
<code>tabstop</code> (<code>ts</code>) specifies how wide tabs characters are
displayed. <code>softtabstop</code> (<code>sts</code>) specifies the
&#8216;virtual tabstop&#8217;, the tabstop that you feel when you press Tab or
Backspace. <code>shiftwidth</code> (<code>sw</code>) specifies the indentation
applied when using autoindent (e.g. the line after an if statement will get
indented by this amount).

For example, while working with Exaile, you have to <code>:set et sts=4
sw=4</code> because Exaile always uses four spaces for indentation. However,
while editing a makefile you would normally set <code>noet</code> and use the
same number for <code>ts</code>, <code>sts</code>, and <code>sw</code>.
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
