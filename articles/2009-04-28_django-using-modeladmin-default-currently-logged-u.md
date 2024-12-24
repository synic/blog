---
title: "Django: Using ModelAdmin to default to currently logged in user"
publishedAt: 2009-04-28T17:51:36-07:00
tags: [Programming, Python, Django]
summary: |
  Yeah, you may have noticed that I've been working on the blog lately.
  Poor openclue.org got flooded with already posted RSS feeds again. This
  happens all to often. Sorry guys.
---
Yeah, you may have noticed that I've been working on the blog lately.
Poor openclue.org got flooded with already posted RSS feeds again. This
happens all to often. Sorry guys.

Anyway, this blog system has the ability to have more than one person
post articles. In the past, sjohannes used to post here too. The model
for articles looks something like this:

```python
from django.contrib.auth.models import User
from django.db import models
import datetime

class Article(models.Model):
    title = models.CharField(max_length=250)
    slug = models.SlugField()
    article = models.TextField()
    pub_date = models.DateTimeField(
        default=datetime.datetime.now)
    publisher = models.ForeignKey(User)
```

As you can see, there's a ForeignKey to
django.contrib.auth.models.User. The question here is: how do I make the
`publisher` field default to the currently logged in user? It was
harder than I thought it should be, but it can be done using hooks in
the ModelAdmin for the Article model. Take a look:

```python
class ArticleAdmin(admin.ModelAdmin):
    list_display = ('pub_date', 'publisher', 'title', 'slug')
    search_fields = ('title', 'article')
    prepopulated_fields = {
        'slug': ('title',),
    }

    def get_form(self, req, obj=None, **kwargs):
        # save the currently logged in user for later
        self.current_user = req.user
        return super(ArticleAdmin, self).get_form(req, obj, **kwargs)

    def formfield_for_dbfield(self, field, **kwargs):
        from django import forms
        from django.contrib.auth import models
        if field.name == 'publisher':
            queryset = models.User.objects.all()
            return forms.ModelChoiceField(
                queryset=queryset, initial=self.current_user.id)
        return super(ArticleAdmin, self).formfield_for_dbfield(field, **kwargs)

admin.site.register(models.Article, ArticleAdmin)
```

As you can see, I overload `get_form` for one purpose: It takes the
HttpRequest as the first argument, which we can use to save the
currently logged in user. The other method that's overloaded,
`formfield_for_dbfield` is called for every field in the model, and
is made for the purpose of specifying your own custom form fields (and
widgets). In this case, we use the same field type that the admin would
have used - `ModelChoiceField` and give it an initial, which is the id
of the currently logged in user.

There you have it. It's a little hacky, but it's the only way I could
find to do it. If you've seen a cleaner way, let me know!

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
