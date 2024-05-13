<!-- :metadata:

title: Amazon AWS API Changes and You
tags: Exaile, Programming
published: 2009-05-08T17:21:47-0700
summary:

I got an email from Amazon today saying that they were changing some things in
their API.  Aside from changing the name of the services from "Amazon Web
Service" to "Product Advertising API", by August 15th, they are requiring that
users of the API send the previously optional HMAC signature when
authenticating.  What does this mean to you?

-->

I got an email from Amazon today saying that they were changing some things in
their API.  Aside from changing the name of the services from "Amazon Web
Service" to "Product Advertising API", by August 15th, they are requiring that
users of the API send the previously optional HMAC signature when
authenticating.  What does this mean to you?

This signature is created using your "Secret Access Key", which is available on
their website in your account details.  From their website: ??"Your Shared
Access Key is secret, and should only be known by you and AWS."??

**This means that in order to use the Product Advertising API (Amazon Web
Services), you will need to sign up for an AWS account.  No more is the
zeroconf album art downloading in programs like Exaile and Amarok**.

Time to move to Last.FM for album art.

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
