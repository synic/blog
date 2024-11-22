<!-- :metadata:

title: Posting articles to Twitter via Django
tags: Programming, Python, Django
publishedAt: 2008-10-12T20:23:51-0700
summary:

I noticed that every time Clint Savage makes a blog update, he posts the URL to
Twitter twice (yeah, that's you herlo :P).  The URLs for each Twitter post are
different, so I figured it must be some sort of automated Wordpress script with
a bug in it...

-->

<p>I noticed that every time Clint Savage makes a blog update, he posts the URL
to Twitter twice (yeah, that's you herlo :P).  The URLs for each Twitter post
are different, so I figured it must be some sort of automated Wordpress script
with a bug in it.</p>
 <p>So, I decided to write something to do the same for
my own blog.  Clint's passes the URL through tinyurl to get an address that's
not too long for Twitter, but in my blog's case, http://vimtips.org/article_id
works just fine.</p>
 <p>Here's the code I used to do it:</p>

```python
from vimtips import settings

try:
    import twitter
except ImportError:
    twitter = None

def twitter_post(article_id, title):
    """
        Posts an article update to twitter
    """

    if not twitter or not hasattr(settings, 'TWITTER_USER') or \
        not hasattr(settings, 'TWITTER_PASS'):
        return

    try:
        api = twitter.Api(settings.TWITTER_USER, settings.TWITTER_PASS)
        api.PostUpdate(
            "Blog Update: %s (http://vimtips.org/%d)" % (title, article_id)
        )
    except Exception, e:
        if settings.DEBUG:
            raise(e)

class Article(models.Model):
    """
        The model representing each news article on the blog
    """

    # ... model fields here ...

    def save(self, *args):
        """
            Saves the article.  If this is a new article, it also posts to
            twitter
        """
        post = not self.id

        models.Model.save(self, *args)

        if post:
            twitter_post(self.id, self.title)
```

<p>This requires that you install the python-twitter module from <a
href='http://code.google.com/p/python-twitter/'>http://code.google.com/p/python-twitter/</a>.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
