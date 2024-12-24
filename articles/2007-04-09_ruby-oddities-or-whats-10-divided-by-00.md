---
title: Ruby oddities (or "what's 1.0 divided by 0.0?")
publishedAt: 2007-04-09T09:37:36-07:00
tags: [Programming]
---
Working very closely with the source code of MRI (Matz's Ruby Implementation, a
term used to distinguish the main Ruby implementation from the language itself)
has exposed me to some of Ruby's most... interesting behaviour. This is one:

```ruby
irb(main):001:0> 1.0 / 0.0
=> Infinity
irb(main):002:0> 1.0.divmod(0.0)[0]
=> NaN
irb(main):003:0> 1.0.div(0.0)
FloatDomainError: Infinity
        from (irb):3:in `div'
        from (irb):3
```

And no, Infinity does not equal `NaN`, and they're not related in terms of
inheritance or anything like that; they're completely different values.
Apparently, `/` follows the IEEE-754 standard (including in cases involving a
negative zero), `divmod` probably follows intuitive mathematics (which I'm
quite sure is not really correct because "1 / 0" has no meaning in strict
mathematics), while @div@ follows a practice quite common in programming
languages of raising an exception (though an odd one at that). Coupled with the
fact that `1 / 0` raises a different exception (ZeroDivisionError), well... go
figure.

<br /><br />
<div class='vimtip'>

<h3>vim tip: <b>Modelines</b></h3>

<p>
In Vim, you can set per-file options using modelines. (This has nothing to do
with the term "mode line" in Emacs. Emacs has a similar feature, but I don't
know what it's called there.) Basically they are special strings in the first
or last few lines of a file that Vim interprets into options. They look like,
for example, <code>// vim: expandtab</code>. The exact syntax is explained in detail in
the [modeline](http://vimdoc.sourceforge.net/htmldoc/options.html#modeline)
section of Vim's help.
</p>

</div>
