---
title: How I made this blog super speedy
slug: how-i-made-this-blog-super-fast
tags: [Programming, Go]
publishedAt: 2026-02-12T12:17:06-07:00
summary: |
  When I was making the current version of this blog, there was a lot of
  discussion on Twitter around improving load times and responsiveness on dynamic
  websites, complete with lots of demonstrations of instant load times even for
  things that access the database, like searches. I decided to see what I could
  pull off with this little site, served from a cheap VPS. I was able to get it
  to a perfect score on Lighthouse, and indeed, if you click around on the links
  or use the search box, you will see that things load pretty much instantly.

  <img src="/static/img/screenshots/lighthouseresults.webp" alt="Lighthouse Results" width="600" height="193" />
---
When I was making the current version of this blog, there was a lot of
discussion on Twitter around improving load times and responsiveness on dynamic
websites, complete with lots of demonstrations of instant load times even for
things that access the database, like searches. I decided to see what I could
pull off with this little site, served from a cheap VPS. I was able to get it
to a perfect score on Lighthouse, and indeed, if you click around on the links
or use the search box, you will see that things load pretty much instantly.

<img src="/static/img/screenshots/lighthouseresults.webp" alt="Lighthouse Results" width="600" height="193" />

## How

I used a few different technologies and ideas:

1. Go backend. As of today, the entire site is served from a single statically
   compiled binary that is only 7MB. Even the assets (images, CSS, JavaScript,
   articles) are embedded directly in the binary.
2. The articles themselves are pre-rendered in the build process, embedded
   directly in the compiled binary, and are all loaded into memory during boot,
   where they stay.

   There is no database access for individual articles or searching, everything
   is accessed directly from memory. I did some tests to make sure this was going
   to be reasonable, even with a ton of articles, and I found that I could
   probably write a couple articles a week for the rest of my life without causing
   any real problems. The storage adapter can be change to use a database if it
   ever becomes absolutely necessary, and then I can just cache frequently accessed
   articles.

   At some point I will add a comment system, and while those
   _will_ be in a database, they will only be loaded on demand when the viewer
   clicks on the "View Comments" button.
3. When you are navigating around the site, I used HTMX to only change the
   parts of the page that are different; the header, general layout, and
   footer do not change.
