---
title: Migrating from Trac to Launchpad
slug: migrating-from-trac-to-launchpad
publishedAt: 2007-09-01T18:57:19-07:00
tags: [Random]
summary: |
  <p>I love trac.  <span class="caps">LOVE</span> it.  With the recent exaile.org
  hack, however, I wanted my bugs and code to be in a place that isn&#8217;t
  going anywhere soon.  I chose Launchpad.</p>
---
<p>I love trac.  <span class="caps">LOVE</span> it.  With the recent exaile.org
hack, however, I wanted my bugs and code to be in a place that isn&#8217;t
going anywhere soon.  I chose Launchpad.</p>

<p>I was faced with a problem &#8211; all of our bugs were in trac.  Lots and
lots of bugs.  I had to somehow migrate them from trac to Launchpad, so I wrote
up a script to do so fairly painlessly.  You use it like this:</p>

```bash
$ ./trac2lp lpusername lppassword project /path/to/trac.db
```

Here&#8217;s the script:</p>

```python
#!/usr/bin/env python
from pysqlite2 import dbapi2 as sqlite
from mechanize import Browser
import sys

try:
    username = sys.argv[1]
    password = sys.argv[2]
    project = sys.argv[3]
    db_location = sys.argv[4]
except:
    print "Useage: trac2lp.py username password project trac.db"
    sys.exit(1)

br = Browser()
br.open("https://launchpad.net/+login")
br.select_form(name="login")
br["loginpage_email"] = username
br["loginpage_password"] = password
response = br.submit()

db = sqlite.connect(db_location)
cur = db.cursor()
cur.execute("SELECT id, summary, reporter, description FROM "
    "ticket WHERE resolution IS NULL ORDER BY id ASC")

# Here, we can keep track of tickets that have already been processed so that
# if something goes wrong, they don’t get processed again
tickets = []
for line in open("tickets.txt").readlines():
    line = line.strip()
    tickets.append(int(line))

h = open("tickets.txt", "a")

for row in cur.fetchall():
    if int(row[0]) in tickets: continue
    h.write("%d" % row[0])

    try:
        br.open("https://launchpad.net/%s/+filebug" % project)
        br.select_form(nr=2)
        br["field.title"] = row[1]
        response = br.submit()

        br.select_form(nr=2)
        br["field.title"] = row[1]
        br["field.comment"] = "%s\n\n\n%s\n%s" % (row[2],
            "This ticket was migrated from the old trac: re #%d" % row[0],
            "Originally reported by: %s" % row[3])

        try:
            br.find_control("field.actions.this_is_my_bug").disabled = True
            control = br.find_control("field.bug_already_reported_as")
            control.items[len(control.items) – 1].selected = True
        except:
            pass
        response = br.submit(id="field.actions.submit_bug")
    except:
        pass
```

<b>Note:</b> This only migrates open tickets.  </p>

<div class='vimtip'>

<h3><b>vim tip:</b> <i>Ron colorscheme</i></h3>

<p>
If you&#8217;re like me, you are used to having your terminal being a white
foreground on a black background.  When using vim in a terminal, I&#8217;ve
found that the default colorscheme is hard on the eyes, or just plain hard to
read with a black background.  I tried out all the schemes that Vim comes with,
and the winner is Ron.  Try it:  <code>:colorscheme ron</code>.  <span
class="caps">IMHO</span>, much better on the eyes.<br />
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
