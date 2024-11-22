<!-- :metadata:

title: Django template autoescaping
tags: Programming, Python, Django
publishedAt: 2007-11-15T23:38:50-0700
summary:


Vimtips.org is running the SVN version of Django.  This morning I ran an
<code>svn update</code>, and I ran into my first API change.  While looking at
my site later on in the day, I noticed that both of my template filters were
being HTML escaped, IE, things like <b>&lt;</b> were showing up as
<b>&amp;lt;</b>...

-->

Vimtips.org is running the SVN version of Django.  This morning I ran an
<code>svn update</code>, and I ran into my first API change.  While looking at
my site later on in the day, I noticed that both of my template filters were
being HTML escaped, IE, things like <b>&lt;</b> were showing up as
<b>&amp;lt;</b>.<br>

My two filters are the pygments highlighting filter (you can see that in action
in <a href='/1'>this</a> article) and the filter that creates the category list
at the end of every article (Under this article, it says "Filed Under:
Programming, Python, Django").<br><br>
 Looking through the svn changelog, I
noticed that they implemented a new feature, called autoescape, which will make
every template variable and custom filters autoescape for safety.  Using:

```jinja
{% autoescape off %}
<a href='{{ link.url }}'>{{ link.name }}</a>
{% endautoescape %}
```

... you can turn off autoescaping.  You can also use the Django template filter
<b>safe</b>.  As for custom filters, to make it so your returned string isn't
autoescaped, you have to mark it as safe.  Here I'm showing my category list
filter with the new `safestring.mark_safe()` function:<br />

```python
from django.utils import safestring

@register.filter(name='category_list')
def category_list(categories):
    """
        Shows all categories as a list of links separated by commas
    """
    c = []
    for category in categories:
        c.append("<a href='/category/%d'>%s</a>" % (category.id,
            category.name))

    return safestring.mark_safe(", ".join(c))
```

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
