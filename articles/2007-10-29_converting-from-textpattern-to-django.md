---
title: Converting from TextPattern to Django
publishedAt: 2007-10-29T23:33:32-07:00
tags: [Random, Django]
summary: |
  I like TextPattern a lot, but it doesn't seem to work well for programmers.  I
  couldn't ever find a syntax highlighting plugin (that worked) for it, and even
  when I did figure out a way to post code TextPattern would try to format it.

  So, I finally had a reason to learn Django, and here is the
  product.  I even implemented my own syntax highlighting filter (<a
  href='http://stderr.ws'>Josh Simpson's</a> idea to do this is actually what
  finally made me want to switch away from TextPattern in the first place):
---
I like TextPattern a lot, but it doesn't seem to work well for programmers.  I
couldn't ever find a syntax highlighting plugin (that worked) for it, and even
when I did figure out a way to post code TextPattern would try to format it.

So, I finally had a reason to learn Django, and here is the
product.  I even implemented my own syntax highlighting filter (<a
href='http://stderr.ws'>Josh Simpson's</a> idea to do this is actually what
finally made me want to switch away from TextPattern in the first place): <br
/><br />

```python
from django import template
from pygments import highlight
from pygments.lexers import get_lexer_by_name
from pygments.formatters import HtmlFormatter
from django.template import Context, loader
from django.template.defaultfilters import stringfilter
import re

register = template.Library()

@register.filter(name='code_highlight')
@stringfilter
def code_highlight(value):
    """
        Checks for <source lang='lang'> tags in the article, and runs them
        through pygments for syntax highlighting
    """
    t = loader.get_template('codeblock.html')

    regex = re.compile(r'(<source lang=([\'"])(\w+)\2>(.*?)</source>)',
        re.DOTALL)
    items = regex.findall(value)

    for (all, crap, lang, text) in items:
        lexer = get_lexer_by_name(lang, stripall=True)
        formatter = HtmlFormatter(linenos=True, cssclass="syntax")
        result = highlight(text, lexer, formatter)

        c = Context({
            'code_block': result
        })

        value = value.replace(all, t.render(c))

    return value
```

There's still a lot of work to do on it, but I'll get the old posts migrated soon.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
