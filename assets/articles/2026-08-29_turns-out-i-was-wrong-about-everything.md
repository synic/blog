---
title: Turns Out, I Was Wrong About Everything
slug: turns-out-i-was-wrong-about-everything
tags: [Programming, Personal, Random]
#publishedAt: 2026-08-29T14:29:56-06:00
---
<!-- summary render-in-body=false -->

I've come to realize that I often have it all wrong. I like to think that I am
pretty open minded, but looking back on my life, there are a ton of things that
I turned my nose up at first, only to come around on them later. From cheese
(the food) to [Docker](https://www.docker.com/) containers, the list of things I
wrote off at first but eventually ended up loving is a mile long.

<!-- end-summary -->

When I was a kid, I was convinced that I hated cheese. And chicken, bananas, and
mushrooms. Obviously, I was wrong, but I didn't figure that out until I was in
my 20s. What happened?

Well, I believe the main thing in these cases is that I just didn't like the
brand my mom bought, or the way she prepared it. For example, she liked bananas
a little on the green side; as soon as they started getting a little brown on
the skin, they were relegated to bread duty or tossed in the trash. It didn't
even occur to me to try them any other way. Wasn't until I was older that I
realized that when the skin starts getting brown spots on it, that's when the
banana becomes sweet and delicoous.

Another thing I hated when I was younger is 80s music. To me, at the time, it
sounded like a game show host trying to sing a song with the most game showy
voice he could possibly muster. Its all my girlfriend wanted to listen to and I
did not get it *at all*. I guess this is just a matter of changing tastes; but I
admit it now, 80s music is pretty fantastic. Don't tell Nancy, I'll never hear
the end of it.

In not quite the same fashion, there are a lot of technologies I've come across
in my career that I intentionally avoided, sometimes for years, only later to
end up putting them in my list of essential software.

A couple examples:

1. [Docker](https://www.docker.com/) - as I mentioned in the intro, I disliked
   Docker at first. I was working at a company that had just gotten their
   [Vagrant](https://github.com/hashicorp/vagrant) setup working quite well for
   spinning up dev environments. Docker had just come out, and someone decided
   to switch us over. One day I came into work, and instead of Vagrant, I now
   had to re-do my environment with Docker. For this reason, I avoided it for a
   while, not fully understanding why I suddenly needed 5 different virtual
   machines when we had just 1 on Vagrant. Obviously I didn't fully understand
   what Docker was or how it worked at that point. Today, I don't know what I
   would do without it.
2. [Tailwind](https://tailwindcss.com/) - this one is pretty obvious. I'm not
   the only one that hated it at first. It grows on you, though.
3. [Go](https://go.dev/) - Oh Go. My first experience with it was during a
   hackathon. On a team of two, we decided to go with Go just because we didn't
   know it and thought it would be fun to learn a new language. This was a poor
   choice. There wasn't enough time to learn it properly, and by the time I was
   done, I was absolutely convinced it was a garbage language. The way the code
   looked to me was hideous. I didn't understand why these fancy structs had
   uppercase variable names and lowercase variable names all mixed together.
   Cant they not even follow their own naming standards? Later on, when I was
   trying to figure out a solution for this blog, I wanted something that would
   still be relevant in several years. This blog is old, and it's gone throught
   several different stacks and only by sheer luck do I still have all the
   articles. For this pass at it, though, I wanted a stable language that didn't
   require a lot of upkeep and package updating. That if I didn't end up writing
   articles for a while, the language wouldn't change so much that I'd need to
   do a major overhaul every 6 months or so just to keep up with security
   patches. Go seemed like it checked all my boxes, so I set out to learning it.
   Once I knew why everything was the way it was (including the case mismatch in
   struct variables names: names that start with a capital letter are public,
   and names that start with a lowercase letter are private to the struct), I
   began to really like it. At this point, I'd probably put it in the list of my
   favorite languages. It's almost like you get C performance out of a language
   that feels a lot like writing Python or some other high level scripting
   language. And the deployment story is unmatched; Go, in most cases, produces
   a binary that has no dependencies. You can even embed all your assets in the
   binary if you want.

The list goes on. So many things along the way could have helped me level up if
I had just ignored that voice that always says *"ew it's a change and I don't
like change"*. I definitely need to find a way to recognize when that type of
thought is happening regarding a new technology, and not write it off before I
really kick the tires.
